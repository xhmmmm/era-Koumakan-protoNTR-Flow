package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	defaultRepoURL = "https://github.com/xhmmmm/era-Koumakan-protoNTR-Flow.git"
	defaultBranch  = "main"
	remoteName     = "origin"
	configFile     = "更新配置.txt"
)

// AppConfig 更新工具配置
type AppConfig struct {
	RepoURL string
	Branch  string
}

func main() {
	fmt.Println("====================================")
	fmt.Println("  era-Koumakan-protoNTR-Flow 更新工具")
	fmt.Println("====================================")
	fmt.Println()

	// 获取程序所在目录（而非当前工作目录）
	exePath, err := os.Executable()
	if err != nil {
		fail("无法确定程序位置", err)
	}
	dir := filepath.Dir(exePath)
	exeName := filepath.Base(exePath)

	// 清理上次运行可能残留的旧文件
	oldExe := exePath + ".old"
	os.Remove(oldExe)

	// 加载可选的配置文件
	cfg := loadConfig(filepath.Join(dir, configFile))
	if cfg.RepoURL == "" {
		cfg.RepoURL = defaultRepoURL
	}
	if cfg.Branch == "" {
		cfg.Branch = defaultBranch
	}

	// 自动检测系统代理：有配置就走代理，没配置就直连
	// 仅设置进程内环境变量，进程退出后即消失，不修改系统环境变量
	if proxy := getSystemProxy(); proxy != "" {
		os.Setenv("HTTP_PROXY", proxy)
		os.Setenv("HTTPS_PROXY", proxy)
		fmt.Printf("检测到系统代理: %s\n", proxy)
	} else {
		fmt.Println("未检测到系统代理，将直连GitHub")
	}

	fmt.Printf("仓库地址: %s\n", cfg.RepoURL)
	fmt.Printf("目标分支: %s\n", cfg.Branch)
	fmt.Println()

	// 检查 .git 目录是否存在
	gitDir := filepath.Join(dir, ".git")
	info, err := os.Stat(gitDir)

	if err == nil && info.IsDir() {
		// ---- .git 存在：尝试打开仓库 ----
		repo, openErr := git.PlainOpen(dir)
		if openErr != nil {
			fmt.Println("Git仓库无法打开，正在重建...")
			if err := rebuildRepo(dir, cfg, exePath, exeName); err != nil {
				fail("重建仓库失败", err)
			}
		} else {
			// 检查远程地址是否正确
			remote, remoteErr := repo.Remote(remoteName)
			if remoteErr != nil || !remoteURLMatches(remote, cfg.RepoURL) {
				fmt.Println("仓库地址配置不正确，正在重建...")
				if err := rebuildRepo(dir, cfg, exePath, exeName); err != nil {
					fail("重建仓库失败", err)
				}
			} else {
				// 远程地址正确，直接拉取更新
				fmt.Println("正在拉取最新更新...")
				if err := pullUpdate(repo, cfg, exePath, exeName); err != nil {
					fail("拉取更新失败", err)
				}
			}
		}
	} else {
		// ---- .git 不存在：初始化新仓库 ----
		fmt.Println("未找到Git仓库，正在初始化...")
		if err := initRepo(dir, cfg, exePath, exeName); err != nil {
			fail("初始化仓库失败", err)
		}
	}

	fmt.Println()
	fmt.Println("[完成] 更新已完成！")
	waitExit()
}

// rebuildRepo 删除现有 .git 并重新初始化
func rebuildRepo(dir string, cfg AppConfig, exePath, exeName string) error {
	gitDir := filepath.Join(dir, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		return fmt.Errorf("删除旧仓库失败: %w", err)
	}
	return initRepo(dir, cfg, exePath, exeName)
}

// initRepo 初始化新仓库、配置远程地址并拉取
func initRepo(dir string, cfg AppConfig, exePath, exeName string) error {
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		return fmt.Errorf("git init 失败: %w", err)
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{cfg.RepoURL},
	})
	if err != nil {
		return fmt.Errorf("添加远程地址失败: %w", err)
	}

	return pullUpdate(repo, cfg, exePath, exeName)
}

// pullUpdate 拉取远程更新并重置工作区
// 使用MergeReset：只更新已跟踪文件，保留未跟踪文件（如存档sav/等）
func pullUpdate(repo *git.Repository, cfg AppConfig, exePath, exeName string) error {
	fmt.Println("正在下载更新...")

	fetchErr := repo.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("+refs/heads/*:refs/remotes/%s/*", remoteName)),
		},
		Progress: os.Stdout,
	})
	if fetchErr != nil && fetchErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("下载失败: %w", fetchErr)
	}

	// 获取远程分支引用 (refs/remotes/origin/main)
	remoteRefName := plumbing.NewRemoteReferenceName(remoteName, cfg.Branch)
	ref, err := repo.Reference(remoteRefName, true)
	if err != nil {
		return fmt.Errorf("找不到远程分支 %s: %w", cfg.Branch, err)
	}

	// 先创建/更新本地分支引用和HEAD，使Reset能正常工作
	localBranch := plumbing.NewHashReference(plumbing.NewBranchReferenceName(cfg.Branch), ref.Hash())
	if err := repo.Storer.SetReference(localBranch); err != nil {
		return fmt.Errorf("更新分支引用失败: %w", err)
	}
	headRef := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName(cfg.Branch))
	if err := repo.Storer.SetReference(headRef); err != nil {
		return fmt.Errorf("设置HEAD失败: %w", err)
	}

	// 获取工作区
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("无法获取工作区: %w", err)
	}

	// Windows下正在运行的exe不能被覆盖，但可以被重命名
	// 先重命名自身，让Reset能写入新版本，旧文件留待下次启动清理
	if exeName != "" {
		oldExe := exePath + ".old"
		os.Remove(oldExe) // 清理可能残留的
		if renameErr := os.Rename(exePath, oldExe); renameErr != nil {
			fmt.Printf("[提示] 无法重命名更新器自身（%v），将跳过更新器文件\n", renameErr)
		}
	}

	// 重置工作区到远程分支的最新提交
	fmt.Println()
	fmt.Println("正在更新文件...")
	err = wt.Reset(&git.ResetOptions{
		Mode:   git.MergeReset,
		Commit: ref.Hash(),
	})
	if err == git.ErrUnstagedChanges {
		// 本地有未暂存的修改，清空索引后重试以绕过检查
		idx, idxErr := repo.Storer.Index()
		if idxErr == nil {
			idx.Entries = nil
			_ = repo.Storer.SetIndex(idx)
		}
		err = wt.Reset(&git.ResetOptions{
			Mode:   git.MergeReset,
			Commit: ref.Hash(),
		})
	}
	if err != nil {
		// Reset失败，尝试恢复exe
		if exeName != "" {
			os.Rename(exePath+".old", exePath)
		}
		return fmt.Errorf("更新文件失败: %w", err)
	}

	// 处理旧版exe
	if exeName != "" {
		if _, statErr := os.Stat(exePath); os.IsNotExist(statErr) {
			// exe不在仓库tree中（如尚未提交），从.old恢复
			os.Rename(exePath+".old", exePath)
		} else {
			// 新版exe已写入，尝试删除旧文件（Windows下可能因占用失败，留待下次启动清理）
			if removeErr := os.Remove(exePath + ".old"); removeErr != nil {
				fmt.Println("[提示] 旧版更新器将在下次启动时自动清理")
			}
		}
	}

	// 显示最新提交信息
	commit, err := repo.CommitObject(ref.Hash())
	if err == nil {
		fmt.Println()
		fmt.Printf("最新提交: %s\n", strings.TrimSpace(commit.Message))
		fmt.Printf("提交时间: %s\n", commit.Author.When.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// remoteURLMatches 检查远程地址是否与预期一致
func remoteURLMatches(remote *git.Remote, expectedURL string) bool {
	for _, u := range remote.Config().URLs {
		if u == expectedURL {
			return true
		}
	}
	return false
}

// loadConfig 加载同目录下的可选配置文件
// 配置文件格式：每行 key=value，# 开头为注释
// 支持的配置项：
//
//	repo_url=https://...   （自定义仓库地址）
//	branch=main            （分支名称）
func loadConfig(path string) AppConfig {
	var cfg AppConfig
	f, err := os.Open(path)
	if err != nil {
		return cfg
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch strings.ToLower(key) {
		case "repo_url":
			cfg.RepoURL = value
		case "branch":
			cfg.Branch = value
		}
	}
	return cfg
}

func fail(msg string, err error) {
	fmt.Println()
	fmt.Printf("[错误] %s\n%v\n", msg, err)
	waitExit()
	os.Exit(1)
}

func waitExit() {
	fmt.Println()
	fmt.Println("按任意键退出...")
	bufio.NewReader(os.Stdin).ReadByte()
}
