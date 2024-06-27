//go:build windows
// +build windows

package process

import (
	"golang.org/x/sys/windows"
)

// windows check pid
const (
	STILL_ACTIVE = 259
)

func isProcessRunning(pid int) (bool, error) {
	// 打开进程
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return false, err
	}
	defer windows.CloseHandle(handle)

	// 获取进程退出代码
	var exitCode uint32
	err = windows.GetExitCodeProcess(handle, &exitCode)
	if err != nil {
		return false, err
	}

	// 如果进程仍然在运行，退出代码为 STILL_ACTIVE (259)
	return exitCode == STILL_ACTIVE, nil
}
