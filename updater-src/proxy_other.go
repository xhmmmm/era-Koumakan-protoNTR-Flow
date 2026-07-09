//go:build !windows

package main

// getSystemProxy 非Windows平台的占位实现
func getSystemProxy() string {
	return ""
}
