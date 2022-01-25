package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

func BuildBinary(wd, repo, branch, dest string) error {
	cache := fmt.Sprintf("%s/%s", wd, DefaultBuildCache)
	err := os.Mkdir(cache, 0755)
	if err != nil {
		return err
	}
	defer os.RemoveAll(cache)

	fmt.Println("Clone the repo from github")
	executable := exec.Command("git", "clone", repo, "-b", branch, cache)
	if err := executable.Start(); err != nil {
		fmt.Println(err)
		return err
	}
	executable.Wait()

	fmt.Println("making the binary")
	executable = exec.Command("make", "build")
	executable.Dir = cache
	if err := executable.Start(); err != nil {
		fmt.Println(err)
		return err
	}
	executable.Wait()

	executable = exec.Command("mv", "build/exchaincli", "build/exchaind", dest)
	executable.Dir = cache
	if err := executable.Start(); err != nil {
		fmt.Println(err)
		return err
	}
	executable.Wait()

	fmt.Println("Build successfully")

	return nil
}

func BuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [dir]",
		Short: "build specific version binary in current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, _ := cmd.Flags().GetString(FlagWorkdir)
			dest := wd
			if len(args) == 1 {
				dest = args[0]
				if !path.IsAbs(dest) {
					dest = fmt.Sprintf("%s/%s", wd, dest)
				}
			}
			repo, _ := cmd.Flags().GetString(FlagGitRepo)
			branch, _ := cmd.Flags().GetString(FlagGitBranch)

			return BuildBinary(wd, repo, branch, dest)
		},
	}

	return cmd
}
