coco-neo4j-backup
=================

Docker Image for automated neo4j backups to S3.

**NB. As of 2016-05-20, this is a work in progress.**

[![CircleCI](https://circleci.com/gh/Financial-Times/coco-neo4j-backup.svg?style=svg)](https://circleci.com/gh/Financial-Times/coco-neo4j-backup)


Requirements
------------

* runs in CoCo
* runs in a different container, but mounts the neo container's filesystem
so that it can do rsync and all that jazz
* uses fleet Go API to shut down neo4j
* needs to tell the deployer not to auto-restart neo4j; there is currently
no API for this.


Running a Backup
----------------

At the time of writing, the neo4j backup process is not automated. It is also tied to the red neo4j instance.

To run a backup (using the `semantic` cluster as an example):

1. SSH to the cluster:

        ssh core@semantic-tunnel-up.ft.com

1. Stop the deployer, red ingester, and red annotator:

        fleetctl stop deployer.service
        fleetctl stop content-ingester-neo4j-red@1.service v1-content-annotator-red@1.service v2-content-annotator-red@1.service

1. Run the backup:

        fleetctl start neo4j-backup.service

1. Wait for it to finish (it should take about half an hour):

        fleetctl journal -f neo4j-backup.service

1. Start the red annotators, the red ingester and the deployer:

        fleetctl start deployer.service
        fleetctl start content-ingester-neo4j-red@1.service v1-content-annotator-red@1.service v2-content-annotator-red@1.service


How to run a restore
--------------------

This is a cold process, so neo4j and its services will be unavailable while this is happening.

1. Warn everybody.
1. SSH into the cluster:

        ssh core@semantic-tunnel-up.ft.com

1. Stop the deployer:

        fleetctl stop deployer.service

1. Identify the machines in the cluster that are running the red and blue instances:
 
        core@ip-172-24-12-247 ~ $ fleetctl list-units | grep ^neo4j-red@
        neo4j-red@1.service                e4882255.../172.24.184.250    active    running

1. Stop the red and blue dependent services:

        fleetctl stop content-ingester-neo4j-red@1.service v1-content-annotator-red@1.service v2-content-annotator-red@1.service
        fleetctl stop content-ingester-neo4j-blue@1.service v1-content-annotator-blue@1.service v2-content-annotator-blue@1.service

1. Stop neo red and blue instances:

        fleetctl stop neo4j-blue@1.service neo4j-red@1.service

1. Download the backup file from S3 (~4GB at the time of writing):

        wget <URL>

1. Repeat these instructions for the red and blue neo4j instances (using red as an example):

    1. SCP the backup file to the cluster host that was running the red neo4j instance.
    1. SSH to the aforementioned red cluster host.
    1. Back up the old data directory and extract the contents of the backup tarball into the `/vol` partition:
    
            /vol/neo4j-red-1
            mv graph.db graph.db.old
            tar -xzvf ~/backup.tar.gz

    1. Start up neo:

            fleetctl start neo4j-red@1.service

    1. Start up dependent services: 

            fleetctl start content-ingester-neo4j-red@1.service v1-content-annotator-red@1.service v2-content-annotator-red@1.service

    1. Test that everything is ok, e.g. by getting some sample data from the red API (you will need to flesh the below curl command
    out with the rest of the parameters required:

            curl -H "<auth header>" https://semantic-up.ft.com/__people-rw-neo4j-red/people/{uuid}

1. Repeat the process for the blue system.
1. Once you're happy that the restore process has worked properly, you can start the deployer and let everyone know that the world
is ok again:

        fleetctl start deployer.service


Development Setup
-----------------

1. Install Go, fleetctl and IntelliJ.
1. Clone this repository.
1. Open the project up in IntelliJ.
1. Set up an SSH tunnel with a dynamic forwarding rule on port 1080.
1. Build and run:

        go build
        ./coco-neo3j-bakup --socksProxy localhost:1080


TODO
----

1. Shut down neo4j's dependencies.
1. Start up neo4j's dependencies.
1. Shameless plagiarise mongo-backup.timer to create neo4j-backup.timer
1. Stop and start the deployer programmatically to avoid neo4j being accidentally started up during a backup.
1. Upload backups into a folder inside the bucket.
1. Write a health check.
1. Lock down the version in services.yaml to a specific tag.
1. Write more tests. Always more tests.


Ideas for automated tests
-------------------------

1. A test that instantiates neo4j, writes some simple data, backs it up, restores it, and tests that it still works as desired.


Notes and Questions
-------------------

1. This thing has to run on the same box as neo4j, right? Is that possible/easy to do in a container-based world?

    * A: Actually, it'll run in its own container and mount the neo4j volume to access the files.
    
2. Why does my IntelliJ build fail when I try to access a function in another file in the same directory? See below:

        $ go build -o "/private/var/folders/rt/8c3952t54cd5q7x08z4m6j5m0000gn/T/Build main.go and rungo" /Users/dafydd/dev/go/src/github.com/Financial-Times/coco-neo4j-backup/main.go
        src/github.com/Financial-Times/coco-neo4j-backup/main.go:24: undefined: backup

    * A: To fix this problem, change the working directory for the run configuration to be the home directory of the project.


Dependencies
------------

This service needs access to the neo4j file system. It therefore relies on the `/vol` partition being present on the host machine,
so that it can be mounted into the container for the `rsync` process. The original plan was 

[fleet-states]: https://github.com/coreos/fleet/blob/master/Documentation/states.md
[docker-hub]: https://hub.docker.com/r/coco/coco-neo4j-backup/
[circle-ci]: https://circleci.com/gh/Financial-Times/coco-neo4j-backup
