//go:build windows

package main

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

// getSystemProxy 读取Windows系统代理设置（注册表）
// 对应"设置 > 网络和 Internet > 代理 > 手动设置代理"
// 返回代理URL（如 http://127.0.0.1:7890），未配置则返回空字符串
func getSystemProxy() string {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()

	// ProxyEnable: 1=启用代理, 0=禁用
	proxyEnable, _, err := k.GetIntegerValue("ProxyEnable")
	if err != nil || proxyEnable == 0 {
		return ""
	}

	// ProxyServer 格式可能是:
	//   "host:port"                    所有协议共用
	//   "http=host:port;https=host:port;ftp=host:port"  分协议
	proxyServer, _, err := k.GetStringValue("ProxyServer")
	if err != nil || proxyServer == "" {
		return ""
	}

	// 处理分协议格式，优先取 https，其次 http
	if strings.Contains(proxyServer, "=") {
		for _, part := range strings.Split(proxyServer, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "https=") {
				return "http://" + strings.TrimPrefix(part, "https=")
			}
			if strings.HasPrefix(part, "http=") {
				return "http://" + strings.TrimPrefix(part, "http=")
			}
		}
		return ""
	}

	// 简单格式 "host:port"
	if !strings.HasPrefix(proxyServer, "http://") && !strings.HasPrefix(proxyServer, "https://") {
		return "http://" + proxyServer
	}
	return proxyServer
}
