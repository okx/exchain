package main

import (
	"os/exec"

	"github.com/spf13/cobra"
)

func KillCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kill",
		Short: "kill all node in this cluster",
		Run: func(cmd *cobra.Command, args []string) {
			path := GetWorkspacePath(cmd, DefaultWorkspace)
			pids := NewPidFile()
			err := pids.Read(path)
			if err != nil {
				cmd.Println("can't not resolve the pid file")
			}
			for _, pid := range pids.Pids {
				excutable := exec.Command("kill", "-9", pid)
				excutable.Run()
			}
		},
	}

	return cmd
}
