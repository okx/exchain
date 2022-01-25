package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type PidFile struct {
	Pids []string
}

func NewPidFile() *PidFile {
	return &PidFile{
		Pids: []string{},
	}
}

func (f *PidFile) Read(workspace string) error {
	path := filepath.Join(workspace, "cluster.pid")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("pid file is not exists")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f.Pids = append(f.Pids, scanner.Text())
	}
	return nil
}

func (f *PidFile) Write(workspace string) error {
	path := filepath.Join(workspace, "cluster.pid")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return errors.New("pid file is exists")
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	for _, line := range f.Pids {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

func (f *PidFile) Delete(workspace string) {
	path := filepath.Join(workspace, "cluster.pid")
	os.Remove(path)
}
