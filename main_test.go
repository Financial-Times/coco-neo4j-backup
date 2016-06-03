package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	log "github.com/Sirupsen/logrus"
)

type mockWriteCloser struct {

}

func (wc *mockWriteCloser) Close() error {
	return nil
}

func (wc *mockWriteCloser) Write(input []byte) (int, error) {
	return len(input), nil
}

func TestEndToEndProcess(t *testing.T) {
	assert := assert.New(t)
	mockFleet := &mockFleetAPI{}
	mockWriter := &mockWriteCloser{}

	rootTestDirName, err := ioutil.TempDir(os.TempDir(), "neoBackupTests")
	log.Infof("Writing to temp directory %s\n", rootTestDirName)
	assert.NoError(err)
	sourceDirPath := filepath.Join(rootTestDirName, "dataFolder")
	targetDirPath := filepath.Join(rootTestDirName, "targetFolder")
	assert.NoError(os.MkdirAll(sourceDirPath, 0777)) // on my Mac this creates directories with 0755 perms. No idea why.
	assert.NoError(os.MkdirAll(targetDirPath, 0777))

	fileContents := "This is a string"
	filename := "test1.txt"
	ioutil.WriteFile(filepath.Join(sourceDirPath, filename), []byte(fileContents), 0666)
	assert.NoError(runInner(mockFleet, mockWriter, sourceDirPath + string(os.PathSeparator), targetDirPath, "mockenv"))
}
