package routes

import (
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/routine"
)

func restartHandler(e *core.RequestEvent) error {
	dev := os.Getenv("DEV")
	restartCmd := os.Getenv("RESTART_CMD")

	if dev == "true" {
		pid := os.Getpid()
		proc, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			return err
		}
	} else {
		if restartCmd == "" {
			return e.InternalServerError("could not find RESTART_CMD env variable", nil)
		}

		e.JSON(http.StatusOK, nil)
		routine.FireAndForget(func() {
			time.Sleep(2 * time.Second)
			args := strings.Split(restartCmd, " ")
			cmd := exec.Command(args[0], args[1:]...)
			e.App.Logger().Info("restarting service")
			cmd.Run()
		})
	}

	return nil
}
