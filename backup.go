package main

import "os/exec"

func rsync() {
	info.Printf("TODO: rsync the neo4j data directory to a temporary location.")

	// TODO Split out the mega-multipack option of "archive" into its carefully selected constituent components.
	cmd := exec.Command("rsync", "--archive", "--verbose", "--delete", "/tmp/foo/", "/tmp/bar")

	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err) // TODO deal with this properly
	}
	info.Printf("Output: %s\n", output)

	info.Printf("TODO: repeat the rsync process until the changes are minimal")
}

func createBackup() {
	info.Printf("TODO NOW DEFINITELY: Create a backup artefact using tar and gzip.")
}

func startNeo() {
	info.Printf("TODO: Start up neo4j.")
	info.Printf("TODO: Start up neo4j's dependencies.")
}

func validateEnvironment() {
	info.Printf("TODO: test that everything is ok")
}