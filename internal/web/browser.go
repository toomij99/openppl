package web

import (
	"fmt"
	"os/exec"
	"runtime"
)

func browserURL(host string, port int) string {
	openHost := host
	if host == "0.0.0.0" || host == "::" || host == "[::]" {
		openHost = "localhost"
	}
	return fmt.Sprintf("http://%s:%d", openHost, port)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform for browser open: %s", runtime.GOOS)
	}

	return cmd.Start()
}
