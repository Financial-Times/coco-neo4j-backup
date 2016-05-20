package main

import (
	"os"
	"time"
	"io"
	"log"
	"flag"
	"github.com/codegangsta/cli"
)

var (
	info *log.Logger
	warn *log.Logger
	fleetEndpoint           = flag.String("fleetEndpoint", "", "Fleet API http endpoint: `http://host:port`")
	socksProxy              = flag.String("socksProxy", "", "address of socks proxy, e.g., 127.0.0.1:9050")
)

const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

func main() {
	initLogs(os.Stdout, os.Stdout, os.Stderr)

	startTime := time.Now()
	info.Printf("Starting backup operation at startTime=%d.\n", startTime)

	// TODO - get awesome code to parse command-line args from @gartonm

	app := cli.NewApp()
	app.Name = "Universal Publishing CoCo neo4j Backup Service"
	app.Usage = "Execute a cold backup of a neo4j instance inside a CoCo cluster and upload it to AWS S3."
	app.Action = func(c *cli.Context) error {
		run()
		return nil
	}

	app.Run(os.Args)
}

func run() {
	rsync()
	shutDownNeo()
	rsync()
	createBackup()
	startNeo()
	uploadToS3()
	validateEnvironment()
	info.Printf("Finishing early because implementation is still on-going.")
}

func uploadToS3() {
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

