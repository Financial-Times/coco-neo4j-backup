package main

import (
	"os"
	"time"
	"io"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"archive/tar"
)

const archiveNameDateFormat = "2006-01-02T15-04-05"

var tarWriter *tar.Writer

func main() {

	startTime := time.Now()
	log.Infof("Starting backup operation at startTime=%d.\n", startTime)

	app := cli.NewApp()
	app.Name = "Universal Publishing CoCo neo4j Backup Service"
	app.Usage = "Execute a cold backup of a neo4j instance inside a CoCo cluster and upload it to AWS S3."
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "fleetEndpoint",
			Value: "http://localhost:49153",
			Usage: "connect to fleet API at `URL`",
			EnvVar: "FLEETCTL_ENDPOINT",
		},
		cli.StringFlag{
			Name: "socksProxy",
			Value: "",
			Usage: "connect to fleet via SOCKS proxy at `PROXY` in IP:PORT format",
			EnvVar: "SOCKS_PROXY",
		},
		cli.StringFlag{
			Name: "awsAccessKey",
			Value: "",
			Usage: "connect to AWS API using access key `KEY`",
			EnvVar: "AWS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name: "awsSecretKey",
			Value: "",
			Usage: "connect to AWS API using secret key `KEY`",
			EnvVar: "AWS_SECRET_KEY",
		},
		cli.StringFlag{
			Name: "dataFolder",
			Value: "/tmp/foo/",
			Usage: "back up from data folder `DATA_FOLDER`",
			EnvVar: "DATA_FOLDER",
		},
		cli.StringFlag{
			Name: "targetFolder",
			Value: "/tmp/bar",
			Usage: "back up to data folder `TARGET_FOLDER`",
			EnvVar: "TARGET_FOLDER",
		},
		cli.StringFlag{
			Name: "s3Domain",
			Value: "",
			Usage: "upload archive to S3 with domain (i.e. hostname) `S3_DOMAIN`",
			EnvVar: "S3_DOMAIN",
		},
		cli.StringFlag{
			Name: "bucketName",
			Value: "",
			Usage: "upload archive to S3 with bucket name `BUCKET_NAME`",
			EnvVar: "BUCKET_NAME",
		},
		cli.StringFlag{
			Name: "env",
			Value: "",
			Usage: "connect to CoCo environment with tag `ENVIRONMENT_TAG`",
			EnvVar: "ENVIRONMENT_TAG",
		},
	}
	app.Action = func(c *cli.Context) error {
		// TODO allow overrides for sourceDir and targetDir
		run(
			startTime,
			c.String("fleetEndpoint"),
			c.String("socksProxy"),
			c.String("awsAccessKey"),
			c.String("awsSecretKey"),
			c.String("dataFolder"),
			c.String("targetFolder"),
			c.String("s3Domain"),
			c.String("bucketName"),
			c.String("env"),
		)
		return nil
	}

	app.Run(os.Args)
}

func run(
	startTime time.Time,
	fleetEndpoint string,
	socksProxy string,
	awsAccessKey string,
	awsSecretKey string,
	dataFolder string,
	targetFolder string,
	s3Domain string,
	bucketName string,
	env string,
	) {

	fleetClient, err := newFleetClient(fleetEndpoint, socksProxy)
	if err != nil {
		panic(err) // TODO handle this properly
	}
	rsync(dataFolder, targetFolder)
	//shutDownNeo(fleetClient)
	rsync(dataFolder, targetFolder)
	// TODO generate archiveName
	archiveName := time.Now().UTC().Format(archiveNameDateFormat)
	archiveName += "_" + env
	archiveName += ".tar.gz"
	pipeReader, _ := createBackup(targetFolder, archiveName)
	startNeo(fleetClient)
	uploadToS3(startTime, awsAccessKey, awsSecretKey, s3Domain, bucketName, archiveName, pipeReader)
	validateEnvironment()
	log.Infof("Finishing early because implementation is still on-going.")
}

func uploadToS3(startTime time.Time, awsAccessKey string, awsSecretKey string, s3Domain string, bucketName string, archiveName string, pipeReader *io.PipeReader) (err error){

	bucketWriterProvider := newS3WriterProvider(awsAccessKey, awsSecretKey, s3Domain, bucketName)

	bucketWriter, err := bucketWriterProvider.getWriter(archiveName)
	if err != nil {
		log.Panic("BucketWriter cannot be created: "+err.Error(), err)
		return err
	}
	defer bucketWriter.Close()

	//upload the archive to the bucket
	_, err = io.Copy(bucketWriter, pipeReader)
	if err != nil {
		log.Panic("Cannot upload archive to S3: "+err.Error(), err)
		return
	}
	pipeReader.Close()

	log.WithFields(log.Fields{
		"archiveName": archiveName,
		"bucketName": bucketName,
		"duration": time.Since(startTime).String(),
	}).Info("Uploaded archive to S3.")
	return nil
}
