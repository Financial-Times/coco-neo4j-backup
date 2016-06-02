package main

import (
	"os/exec"
	"strings"
	"os"
	"io"
	"compress/gzip"
	"archive/tar"
	"path/filepath"
	log "github.com/Sirupsen/logrus"
	"time"
)

func rsync(sourceDir string, targetDir string) (error) {
	startTime := time.Now()
	if ! strings.HasSuffix(sourceDir, "/") {
		log.WithFields(log.Fields{"sourceDir": sourceDir}).Warn("Source directory should probably have a trailing slash!")
	}
	// TODO Split out the mega-multipack option of "archive" into its carefully selected constituent components.
	cmd := exec.Command("rsync", "--archive", "--verbose", "--delete", sourceDir, targetDir)

	output, err := cmd.CombinedOutput()
	o := string(output[:len(output)])
	if err != nil {
		log.WithFields(log.Fields{
			"sourceDir": sourceDir,
			"targetDir": targetDir,
			"output": o,
			"err": err,
		}).Panic("Error executing rsync command!")
	} else {
		log.WithFields(log.Fields{"output": o, "duration_s": time.Since(startTime).String()}).Info("rsync process complete.")
	}
	return err
}

func createBackup(dataFolder string, archiveName string) (*io.PipeReader, error) {
	startTime := time.Now()
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"dataFolder": dataFolder,
			"err": err,
		}).Panic("Directory does not exist!")
		return nil, err
	}
	if _, err := os.Stat(archiveName); os.IsExist(err) {
		log.WithFields(log.Fields{
			"archiveName": archiveName,
			"err": err,
		}).Panic("Archive file already exists!")
		return nil, err
	}

	log.WithFields(log.Fields{"archiveName": archiveName,}).Info("Asynchronously compressing archive.")

	pipeReader, pipeWriter := io.Pipe()
	//compress the tar archive
	gzipWriter := gzip.NewWriter(pipeWriter)
	//create a tar archive
	tarWriter = tar.NewWriter(gzipWriter)

	//a goroutine is needed because the pipe is synchronous:
	//the writer will block until the reader is reading and vice-versa
	go func() {
		//we have to close these here so that the read function doesn't block
		defer pipeWriter.Close()
		defer gzipWriter.Close()
		defer tarWriter.Close()

		//recursively walk the filetree of the data folder,
		//writing all files and folder structure to the archive
		filepath.Walk(dataFolder, addtoArchive)
		log.WithFields(log.Fields{"duration_s": time.Since(startTime).String()}).Info("tar/gzip process complete.")
	}()
	return pipeReader, nil
}

func addtoArchive(path string, fileInfo os.FileInfo, err error) error {
	if fileInfo.IsDir() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		log.WithFields(log.Fields{"path": path, "err": err}).Panic("Cannot open file to add to archive.")
		return err
	}
	defer file.Close()

	//create and write tar-specific file header
	fileInfoHeader, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		log.WithFields(log.Fields{"path": path, "err": err}).Panic("Cannot create tar header.")
		return err
	}
	//replace file name with full path to preserve file structure in the archive
	fileInfoHeader.Name = path
	err = tarWriter.WriteHeader(fileInfoHeader)
	if err != nil {
		log.WithFields(log.Fields{"path": path, "err": err}).Panic("Cannot create tar header.")
		return err
	}

	//add file to the archive
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		log.WithFields(log.Fields{"path": path, "err": err}).Panic("Cannot add file to archive.")
		return err
	}

	log.WithFields(log.Fields{"path": path}).Info("Added file to archive.")
	return nil
}

func validateEnvironment() {
	log.Info("TODO: test that everything is ok")
}