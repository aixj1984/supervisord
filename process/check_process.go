//go:build !windows
// +build !windows

package process

func isProcessRunning(pid int) (bool, error) {
	return true, nil
}
