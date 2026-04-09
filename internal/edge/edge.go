package edge

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Kill termina todos os processos do Microsoft Edge em execução.
func Kill() error {
	// Kill Edge processes multiple times — Edge can respawn background processes
	for range 3 {
		cmd := exec.Command("taskkill", "/F", "/IM", "msedge.exe")
		cmd.Run() // ignore error, process may not exist
		time.Sleep(time.Second * 2)
	}

	err := deleteLastSessions()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 2)

	return nil
}

// deleteLastSessions remove os arquivos de última sessão do Microsoft Edge
func deleteLastSessions() error {
	userDataDir := filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		"Microsoft",
		"Edge",
		"User Data",
	)

	os.Remove(filepath.Join(userDataDir, "SingletonLock"))
	os.Remove(filepath.Join(userDataDir, "SingletonCookie"))
	os.Remove(filepath.Join(userDataDir, "SingletonSocket"))

	err := os.RemoveAll(filepath.Join(userDataDir, "Default", "Sessions"))
	return err
}
