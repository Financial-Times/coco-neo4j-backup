package main

import (
	"archive/tar"
	"flag"
	"io"
	"log"
	"os"
	"time"
	//"github.com/Redundancy/go-sync"
)

var tarWriter *tar.Writer
var info *log.Logger
var warn *log.Logger

const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

func main() {
	initLogs(os.Stdout, os.Stdout, os.Stderr)

	startTime := time.Now()
	info.Printf("Starting backup operation at startTime=%d.\n", startTime)

	awsAccessKey, awsSecretKey, bucketName, dataFolder, s3Domain, env := readArgs()
	printArgs(awsAccessKey, awsSecretKey, bucketName, dataFolder, s3Domain, env)
	checkIfArgsAreEmpty(awsAccessKey, awsSecretKey, bucketName, dataFolder, s3Domain, env)

	//fs := &gosync.BasicSummary{

}

func readArgs() (string, string, string, string, string, string) {
	awsAccessKey := flag.String("awsAccessKey", "", "AWS access key")
	awsSecretKey := flag.String("awsSecretKey", "", "AWS secret key")
	bucketName := flag.String("bucketName", "", "Bucket name")
	dataFolder := flag.String("dataFolder", "", "Data folder to back up")
	s3Domain := flag.String("s3Domain", "", "The S3 domain")
	env := flag.String("env", "", "The environment")
	flag.Parse()
	return *awsAccessKey, *awsSecretKey, *bucketName, *dataFolder, *s3Domain, *env
}

func printArgs(awsAccessKey string, awsSecretKey string, bucketName string, dataFolder string, s3Domain string, env string) {
	info.Println("Using arguments:")
	info.Println("bucketName   : ", bucketName)
	info.Println("dataFolder   : ", dataFolder)
	info.Println("s3Domain     : ", s3Domain)
	info.Println("env          : ", env)
}

func abortOnInvalidParams(paramNames []string) {
	for _, paramName := range paramNames {
		warn.Println(paramName + " is missing or invalid!")
	}
	log.Panic("Aborting backup operation!")
}

func checkIfArgsAreEmpty(awsAccessKey string, awsSecretKey string, bucketName string, dataFolder string, s3Domain string, env string) {
	var invalidArgs []string

	if len(awsAccessKey) < 1 {
		invalidArgs = append(invalidArgs, "awsAccessKey")
	}
	if len(awsSecretKey) < 1 {
		invalidArgs = append(invalidArgs, "awsSecretKey")
	}
	if len(bucketName) < 1 {
		invalidArgs = append(invalidArgs, "bucketName")
	}
	if len(dataFolder) < 1 {
		invalidArgs = append(invalidArgs, "dataFolder")
	}
	if len(s3Domain) < 1 {
		invalidArgs = append(invalidArgs, "s3Domain")
	}
	if len(env) < 1 {
		invalidArgs = append(invalidArgs, "env")
	}

	if len(invalidArgs) > 0 {
		abortOnInvalidParams(invalidArgs)
	}
}

func addtoArchive(path string, fileInfo os.FileInfo, err error) error {
	if fileInfo.IsDir() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		log.Panic("Cannot open file to add to archive: " + path + ", error: " + err.Error(), err)
	}
	defer file.Close()

	//create and write tar-specific file header
	fileInfoHeader, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		log.Panic("Cannot create tar header, error: " + err.Error(), err)
	}
	//replace file name with full path to preserve file structure in the archive
	fileInfoHeader.Name = path
	err = tarWriter.WriteHeader(fileInfoHeader)
	if err != nil {
		log.Panic("Cannot write tar header, error: " + err.Error(), err)
	}

	//add file to the archive
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		log.Panic("Cannot add file to archive, error: " + err.Error(), err)
	}

	info.Println("Added file " + path + " to archive.")
	return nil
}

func initLogs(infoHandle io.Writer, warnHandle io.Writer, panicHandle io.Writer) {
	//to be used for INFO-level logging: info.Println("foor is now bar")
	info = log.New(infoHandle, "INFO  - ", logPattern)
	//to be used for WARN-level logging: info.Println("foor is now bar")
	warn = log.New(warnHandle, "WARN  - ", logPattern)

	//to be used for panics: log.Panic("foo is on fire")
	//log.Panic() = log.Printf + panic()
	log.SetFlags(logPattern)
	log.SetPrefix("ERROR - ")
	log.SetOutput(panicHandle)
}

