package main

import (
	"os/exec"
	"github.com/coreos/fleet/client"
	"strings"
	"os"
	"io"
	"compress/gzip"
	"archive/tar"
	"path/filepath"
	"log"
)

func rsync(sourceDir string, targetDir string) {
	if ! strings.HasSuffix(sourceDir, "/") {
		warn.Printf("Source directory should probably have a trailing slash! sourceDir=\"%s\"", sourceDir)
	}
	info.Printf("TODO: rsync the neo4j data directory to a temporary location.")

	// TODO Split out the mega-multipack option of "archive" into its carefully selected constituent components.
	cmd := exec.Command("rsync", "--archive", "--verbose", "--delete", sourceDir, targetDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err) // TODO deal with this properly
	}
	info.Printf("Output: %s\n", output)

	info.Printf("TODO: repeat the rsync process until the changes are minimal")
}

func createBackup(dataFolder string, archiveName string) {
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		warn.Printf("Directory dataFolder=\"%s\" does not exist!", dataFolder)
		panic(err) // TODO Handle this properly.
	}
	if _, err := os.Stat(archiveName); os.IsExist(err) {
		warn.Printf("Archive file archiveName=\"%s\" already exists!", archiveName)
		panic(err) // TODO Handle this properly.
	}
	info.Printf("TODO NOW DEFINITELY: Create a backup artefact using tar and gzip.")

	_, pipeWriter := io.Pipe()
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
	}()


}

func addtoArchive(path string, fileInfo os.FileInfo, err error) error {
	if fileInfo.IsDir() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		log.Panic("Cannot open file to add to archive: "+path+", error: "+err.Error(), err)
	}
	defer file.Close()

	//create and write tar-specific file header
	fileInfoHeader, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		log.Panic("Cannot create tar header, error: "+err.Error(), err)
	}
	//replace file name with full path to preserve file structure in the archive
	fileInfoHeader.Name = path
	err = tarWriter.WriteHeader(fileInfoHeader)
	if err != nil {
		log.Panic("Cannot write tar header, error: "+err.Error(), err)
	}

	//add file to the archive
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		log.Panic("Cannot add file to archive, error: "+err.Error(), err)
	}

	info.Println("Added file " + path + " to archive.")
	return nil
}

func startNeo(fleetClient client.API) {
	info.Printf("TODO: Start up neo4j.")
	info.Printf("TODO: Start up neo4j's dependencies.")
	// TODO figure out the correct values for these.
	fleetClient.SetUnitTargetState("neo", "active")
}

func validateEnvironment() {
	info.Printf("TODO: test that everything is ok")
}