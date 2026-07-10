package edge

import (
	"os"
	"os/exec"
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

	// err := deleteLastSessions()
	// if err != nil {
	// 	return err
	// }

	time.Sleep(time.Second * 2)

	return nil
}

// CopyUserData cria uma cópia de User Data do MSEdge exclusiva para automação.
func CopyUserData(src string, dst string) error {
	cmd := exec.Command(
		"robocopy",
		src,
		dst,
		"/E",
		"/COPY:DAT",
		"/R:0",
		"/W:0",
		"/XF",
		"lockfile",
		"SingletonLock",
		"SingletonSocket",
		"/NFL",
		"/NDL",
	)

	err := cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			if exitErr.ExitCode() <= 7 {
				return nil
			}
			return err
		}
		return err
	}

	return nil
}

func RemoveTempUserData(src string) error {
	return os.RemoveAll(src)
}

// deleteLastSessions remove os arquivos de última sessão do Microsoft Edge
// func deleteLastSessions() error {
// 	userDataDir := filepath.Join(
// 		os.Getenv("LOCALAPPDATA"),
// 		"Microsoft",
// 		"Edge",
// 		"User Data",
// 	)

// 	os.Remove(filepath.Join(userDataDir, "SingletonLock"))
// 	os.Remove(filepath.Join(userDataDir, "SingletonCookie"))
// 	os.Remove(filepath.Join(userDataDir, "SingletonSocket"))

// 	err := os.RemoveAll(filepath.Join(userDataDir, "Default", "Sessions"))
// 	return err
// }
