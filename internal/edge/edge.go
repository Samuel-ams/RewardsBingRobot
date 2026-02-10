package edge

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Kill termina todos os processos do Microsoft Edge em execução.
func Kill() error {
	cmd := exec.Command("taskkill",
		"/F",
		"/IM",
		"msedge.exe",
	)
	err := cmd.Run()
	if err != nil && err.(*exec.ExitError).ExitCode() != 128 {
		return err
	}

	time.Sleep(time.Second)

	err = deleteLastSessions()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 10)

	return nil
}

// deleteLastSessions remove os arquivos de última sessão do Microsoft Edge
func deleteLastSessions() error {
	err := os.RemoveAll(filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		"Microsoft",
		"Edge",
		"User Data",
		"Default",
		"Sessions",
	))
	return err
}
