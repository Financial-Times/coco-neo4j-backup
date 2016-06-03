package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	log "github.com/Sirupsen/logrus"
	"bytes"
	"compress/gzip"
	"archive/tar"
	"strings"
)

const rootDirName = "neoBackupTests"
const sourceDirName = "dataFolder"
const targetDirName = "targetFolder"
const mockFile1Name = "test1.txt"
const mockFile1Contents = "This is a string"

type mockWriteCloser struct {
	buffer bytes.Buffer
}

func (wc *mockWriteCloser) Close() error {
	return nil
}

func (wc *mockWriteCloser) Write(input []byte) (int, error) {
	return wc.buffer.Write(input)
}

func makeTempFilesystem(rootDirName string) (rootPath string, err error) {
	rootTestDirName, err := ioutil.TempDir(os.TempDir(), rootDirName)
	if err != nil {
		return "", err
	}
	log.Infof("Writing to temp directory %s\n", rootTestDirName)
	sourceDirPath := filepath.Join(rootTestDirName, sourceDirName)
	//targetDirPath := filepath.Join(rootTestDirName, targetDirName)
	err = os.MkdirAll(sourceDirPath, 0777) // on my Mac this creates directories with 0755 perms. No idea why.
	if err != nil {
		return "", err
	}
	//err = os.MkdirAll(targetDirPath, 0777)
	//if err != nil {
	//	return "", err
	//}
	ioutil.WriteFile(filepath.Join(sourceDirPath, mockFile1Name), []byte(mockFile1Contents), 0666)
	return rootTestDirName, nil
}

func TestEndToEndProcessHappyPath(t *testing.T) {
	log.Info("TestEndToEndProcessHappyPath")
	assert := assert.New(t)
	mockFleet := &mockFleetApi{}
	mockWriter := &mockWriteCloser{}

	rootPath, err := makeTempFilesystem(rootDirName)
	assert.NoError(err)

	assert.NoError(runInner(
		mockFleet,
		mockWriter,
		filepath.Join(rootPath, sourceDirName) + string(os.PathSeparator),
		filepath.Join(rootPath, targetDirName),
		"mockenv"))
	assert.NotZero(mockWriter.buffer.Len())

	gzipReader, err := gzip.NewReader(&mockWriter.buffer)
	assert.NoError(err)
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			break
		}
		buf, err := ioutil.ReadAll(tarReader)
		if err != nil {
			t.Fatal("There was a problem reading from the tarball")
		}
		switch {
			case strings.HasSuffix(header.Name, mockFile1Name):
				assert.Equal(buf, []byte(mockFile1Contents))
			default:
				t.Fatalf("Unknown entry encountered inside tarball: %s", header.Name)
		}
	}
}

func TestRunInnerWithNonExistentFolders(t *testing.T) {
	log.Info("TestRunInnerWithNonExistentFolders")
	assert := assert.New(t)
	mockFleet := &mockFleetApi{}
	mockWriter := &mockWriteCloser{}
	assert.NotNil(runInner(
		mockFleet,
		mockWriter,
		"/tmp/doesnotexist/",
		"/tmp/doesnotexisteither",
		"mockenv"))
}

func TestRunInnerWithFleetError(t *testing.T) {
	log.Info("TestRunInnerWithFleetError")
	assert := assert.New(t)
	mockFleetWithErrors := &mockFleetApiError{}
	mockWriter := &mockWriteCloser{}
	rootPath, err := makeTempFilesystem(rootDirName)
	assert.NoError(err)
	assert.EqualError(runInner(
		mockFleetWithErrors,
		mockWriter,
		filepath.Join(rootPath, sourceDirName) + string(os.PathSeparator),
		filepath.Join(rootPath, targetDirName),
		"mockenv"), UnitStateErrorText)
}

func TestRunInnerWithNoTrailingSlash(t *testing.T) {
	log.Info("TestRunInnerWithNoTrailingSlash")
	assert := assert.New(t)
	mockFleet := &mockFleetApi{}
	mockWriter := &mockWriteCloser{}
	rootPath, err := makeTempFilesystem(rootDirName)
	assert.NoError(err)
	assert.NoError(runInner(
		mockFleet,
		mockWriter,
		filepath.Join(rootPath, sourceDirName),
		filepath.Join(rootPath, targetDirName),
		"mockenv"))
}
