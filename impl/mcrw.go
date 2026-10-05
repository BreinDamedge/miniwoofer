package impl

// microwoofer tests

import (
	"fmt"
	"os/exec" // now this can't be safe
	"runtime"
)

func OpenInBrowser(file_path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", file_path)
	case "linux":
		cmd = exec.Command("xdg-open", file_path)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
