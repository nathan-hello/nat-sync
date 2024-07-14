package client

import (
	"fmt"
	"os/exec"
)

type LaunchTargets string

const (
	LaunchMpv LaunchTargets = "mpv"
	LaunchVlc LaunchTargets = "vlc"
)

type Launch struct {
	Player LaunchTargets
}

func launchPlayer(p LaunchTargets, init chan bool) error {
	ipcSocketPath := "/tmp/nat-sync-mpv-socket"

	var cmd *exec.Cmd
	if p == "mpv" {
		cmd = exec.Command("mpv", "--idle", "--force-window", "--input-ipc-server="+ipcSocketPath)
	}

	err := cmd.Start()
	if err != nil {
		return err
	}
	fmt.Println(p, "started")
	init <- true

	cmd.Wait()
	return nil

}
