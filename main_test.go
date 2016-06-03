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

type mockWriteCloser struct {
	buffer bytes.Buffer
}

func (wc *mockWriteCloser) Close() error {
	return nil
}

func (wc *mockWriteCloser) Write(input []byte) (int, error) {
	return wc.buffer.Write(input)
}

func TestEndToEndProcessHappyPath(t *testing.T) {
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
			case strings.HasSuffix(header.Name, "test1.txt"):
				assert.Equal(buf, []byte(fileContents))
			default:
				t.Fatalf("Unknown entry encountered inside tarball: %s", header.Name)
		}
	}
}

// TODO test error cases to improve the coverage