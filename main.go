package main

import (
	"os"
	"time"
	"io"
	"log"
	"github.com/urfave/cli"
)

var (
	info *log.Logger
	warn *log.Logger
)

const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

func main() {
	initLogs(os.Stdout, os.Stdout, os.Stderr)

	startTime := time.Now()
	info.Printf("Starting backup operation at startTime=%d.\n", startTime)

	app := cli.NewApp()
	app.Name = "Universal Publishing CoCo neo4j Backup Service"
	app.Usage = "Execute a cold backup of a neo4j instance inside a CoCo cluster and upload it to AWS S3."
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "fleetEndpoint, f",
			Value: "http://localhost:49153",
			Usage: "connect to fleet API at `URL`",
			EnvVar: "FLEETCTL_ENDPOINT",
		},
		cli.StringFlag{
			Name: "socksProxy, p",
			Value: "",
			Usage: "connect to fleet via SOCKS proxy at `PROXY` in IP:PORT format",
			EnvVar: "SOCKS_PROXY",
		},
		cli.StringFlag{
			Name: "awsAccessKey, a",
			Value: "",
			Usage: "connect to AWS API using access key `KEY`",
			EnvVar: "AWS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name: "awsSecretKey, s",
			Value: "",
			Usage: "connect to AWS API using secret key `KEY`",
			EnvVar: "AWS_SECRET_KEY",
		},
		cli.StringFlag{
			Name: "s3bucket, b",
			Value: "coco-neo4j-backups",
			Usage: "upload archive to S3 bucket `BUCKET`",
			EnvVar: "NEO4J_BACKUP_S3BUCKET",
		},
	}
	app.Action = func(c *cli.Context) error {
		run(
			c.String("fleetEndpoint"),
			c.String("socksProxy"),
			c.String("awsAccessKey"),
			c.String("awsSecretKey"),
			c.String("s3bucket"),
		)
		return nil
	}

	app.Run(os.Args)
}

func run(fleetEndpoint string, socksProxy string, awsAccessKey string, awsSecretKey string, s3bucket string) {
	fleetClient, err := newFleetClient(fleetEndpoint, socksProxy)
	if err != nil {
		panic(err) // TODO handle this properly
	}
	rsync()
	shutDownNeo(fleetClient)
	rsync()
	createBackup()
	startNeo(fleetClient)
	uploadToS3(awsAccessKey, awsSecretKey, s3bucket)
	validateEnvironment()
	info.Printf("Finishing early because implementation is still on-going.")
}

func uploadToS3(awsAccessKey string, awsSecretKey string, s3bucket string) {
	info.Printf("TODO NOW DEFINITELY: Upload the archive to S3.")
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

