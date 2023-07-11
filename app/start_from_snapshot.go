package app

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/klauspost/pgzip"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func prepareSnapshotDataIfNeed(snapshotURL string, home string, logger log.Logger) {
	if snapshotURL == "" {
		return
	}

	snapshotHome := filepath.Join(home, ".download_snapshots")

	// check whether the snapshot file has been downloaded
	byteData, err := os.ReadFile(filepath.Join(snapshotHome, ".record"))
	if err == nil && strings.Contains(string(byteData), snapshotURL) {
		return
	}

	if _, err := url.Parse(snapshotURL); err != nil {
		panic(errors.Wrap(err, "invalid snapshot URL"))
	}

	// download snapshot
	snapshotFile, err := downloadSnapshot(snapshotURL, snapshotHome, logger)
	if err != nil {
		panic(err)
	}

	// uncompress snapshot
	logger.Info("start to uncompress snapshot")
	if err := extractTarGz(snapshotFile, snapshotHome); err != nil {
		panic(err)
	}

	// delete damaged data
	logger.Info("start to delete damaged data")
	if err := os.RemoveAll(filepath.Join(home, "data")); err != nil {
		panic(err)
	}

	// move snapshot data
	logger.Info("start to move snapshot data")
	if err := moveDir(filepath.Join(snapshotHome, "data"), filepath.Join(home, "data")); err != nil {
		panic(err)
	}

	os.Remove(snapshotFile)

	os.WriteFile(filepath.Join(snapshotHome, ".record"), []byte(snapshotURL+"\n"), 0644)

	logger.Info("snapshot data is ready, start node soon!")
}

func downloadSnapshot(url, outputPath string, logger log.Logger) (string, error) {
	// create dir
	_, err := os.Stat(outputPath)
	if err != nil {
		os.MkdirAll(outputPath, 0755)
	}

	fileName := url[strings.LastIndex(url, "/")+1:]
	targetFile := filepath.Join(outputPath, fileName)

	// check file exists
	if _, err := os.Stat(targetFile); err == nil {
		os.Remove(targetFile)
	}

	var stdoutProcessStatus bytes.Buffer

	axel := exec.Command("axel", "-n", fmt.Sprintf("%d", runtime.NumCPU()), "-o", targetFile, url)
	axel.Stdout = io.MultiWriter(ioutil.Discard, &stdoutProcessStatus)
	done := make(chan struct{})
	defer close(done)

	// print download detail
	go func() {
		tick := time.NewTicker(time.Millisecond * 50)
		defer tick.Stop()
		for {
			select {
			case <-done:
				return
			case <-tick.C:
				bts := make([]byte, stdoutProcessStatus.Len())
				stdoutProcessStatus.Read(bts)
				logger.Info(string(bts))
			}
		}
	}()

	// run and wait
	err = axel.Run()
	if err != nil {
		return "", err
	}

	return targetFile, nil
}

func extractTarGz(tarGzFile, destinationDir string) error {
	// open .tar.gz
	file, err := os.Open(tarGzFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// use gzip.Reader
	gzReader, err := pgzip.NewReaderN(file, 1<<22, runtime.NumCPU())
	if err != nil {
		return err
	}
	defer gzReader.Close()

	// create tar.Reader
	tarReader := tar.NewReader(gzReader)

	// uncompress file back to back
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header == nil {
			continue
		}
		target := filepath.Join(destinationDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(target, 0755)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			parent := filepath.Dir(target)
			err = os.MkdirAll(parent, 0755)
			if err != nil {
				return err
			}

			file, err := os.Create(target)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func moveDir(sourceDir, destinationDir string) error {
	sourceInfo, err := os.Stat(sourceDir)
	if err != nil {
		return err
	}

	if !sourceInfo.IsDir() {
		return fmt.Errorf("%s isn't dir", sourceDir)
	}

	_, err = os.Stat(destinationDir)
	if err == nil {
		return fmt.Errorf("dest dir %s exists", destinationDir)
	}

	// move
	err = os.Rename(sourceDir, destinationDir)
	if err != nil {
		return err
	}

	return nil
}
