package main

import (
	"os"
	"time"
	"io"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"archive/tar"
	"fmt"
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
		log.WithFields(log.Fields{
			"fleetEndpoint": fleetEndpoint,
			"socksProxy": socksProxy,
			"err": err,
		}).Panic("Error instantiating fleet client; backup process failed.")
		os.Exit(1)
	}
	err = rsync(dataFolder, targetFolder)
	if err != nil {
		log.WithFields(log.Fields{
			"dataFolder": dataFolder,
			"targetFolder": targetFolder,
			"err": err,
		}).Panic("Error synchronising neo4j files while database is running; backup process failed.")
		os.Exit(1)
	}
	err = shutDownNeo(fleetClient)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Panic("Error shutting down neo4j; backup process failed.")
		os.Exit(1)
	}
	err = rsync(dataFolder, targetFolder)
	if err != nil {
		log.WithFields(log.Fields{
			"fleetEndpoint": fleetEndpoint,
			"socksProxy": socksProxy,
			"err": err,
		}).Panic("Error synchronising neo4j files while database is stopped; backup process failed.")
		os.Exit(1)
	}
	err = startNeo(fleetClient)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Panic("Error starting up neo4j.")
		os.Exit(1)
	}
	archiveName := fmt.Sprintf("neo4j_backup_%s_%s.tar.gz", time.Now().UTC().Format(archiveNameDateFormat), env)
	pipeReader, err := createBackup(targetFolder, archiveName)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Panic("Error creating backup tarball.")
		os.Exit(1)
	}
	log.WithFields(log.Fields{
		"s3Domain": s3Domain,
		"bucketName": bucketName,
		"archiveName": archiveName,
		"err": err,
	}).Panic("Uploading archive to S3.")
	err = uploadToS3(awsAccessKey, awsSecretKey, s3Domain, bucketName, archiveName, pipeReader)
	if err != nil {
		log.WithFields(log.Fields{
			"s3Domain": s3Domain,
			"bucketName": bucketName,
			"archiveName": archiveName,
			"err": err,
		}).Panic("Error uploading to S3; backup process failed.")
		os.Exit(1)
	}
	validateEnvironment()
	log.WithFields(log.Fields{
		"archiveName": archiveName,
		"bucketName": bucketName,
		"duration": time.Since(startTime).String(),
	}).Info("Artefact successfully uploaded to S3; backup process complete.")
}

func uploadToS3(awsAccessKey string, awsSecretKey string, s3Domain string, bucketName string, archiveName string, pipeReader *io.PipeReader) (err error){
	startTime := time.Now()
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
		return err
	}
	pipeReader.Close()

	log.WithFields(log.Fields{
		"archiveName": archiveName,
		"bucketName": bucketName,
		"duration": time.Since(startTime).String(),
	}).Info("Uploaded archive to S3.")
	return nil
}
