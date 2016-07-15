coco-neo4j-backup
=================

Docker Image for automated neo4j backups to S3.

**NB. As of 2016-05-20, this is a work in progress.**

[![CircleCI](https://circleci.com/gh/Financial-Times/coco-neo4j-backup.svg?style=svg)](https://circleci.com/gh/Financial-Times/coco-neo4j-backup)


Documentation
-------------

* [Backing up and restoring](backup-and-restore.md) - instructions for running a backup, and running a restore.


Development Setup
-----------------

1. Install Go, fleetctl and IntelliJ.
1. Clone this repository.
1. Open the project up in IntelliJ.
1. Set up an SSH tunnel with a dynamic forwarding rule on port 1080.
1. Build and run:

        go build
        ./coco-neo3j-bakup --socksProxy localhost:1080

1. Testing that everything builds ok:

        docker build -t $(basename $PWD) .

1. Releasing a new version:

    1. Tag the release according to semantic versioning principles:

            git tag 0.x.0
            git push --tag

    1. Check that Docker Hub built it ok: https://hub.docker.com/r/coco/coco-neo4j-backup/builds/
    1. Update the version in `services.yaml` via a branch/PR.
    1. Wait for the deployer to deploy the service.


TODO
----

The below items may want to be implemented at some point, perhaps when we start "backup 2.0" if/when we start using hot backups
with Neo4j Enterprise.

1. Shamelessly plagiarise `mongo-backup.timer` to create `neo4j-backup.timer`, for scheduled backups.
1. Upload backups into a folder inside the bucket in a format something like `neo4j-<cluster>`, e.g. `neo4j-pre-prod`.
1. Write a health check.
1. ~~Lock down the version in services.yaml to a specific tag.~~ DONE
1. Write more tests. Always more tests.
1. Print a link to the backup archive in S3.
1. Check CPU usage, then see if using an LZ4 compressor reduces CPU usage (potentially at the cost of a larger backup file).
1. Switch to using a library like [env-decode] for much simpler parsing of environment variables without needing CLI params,
which are unnecessary for most apps.
1. Make it possible to back up red *or* blue (rather than just red). (Requested by Scott 2016-07-15)

### TODO items that will probably no longer be necessary once we have hot neo4j backups

1. Add `ionice` in front of the `nice rsync` statement, to further reduce resource usage
(suggested by [martingartonft](https://github.com/martingartonft) on 2016-07-11). NB. will no longer be necessary once we are doing
hot backups, although we might want to run the entire service under a low process priority.
1. Stop and start the deployer programmatically to avoid neo4j being accidentally started up during a backup.
It would be wise to restart `deployer.service` programmatically immediately after neo is started back up, rather than doing it manually
after the backup process is complete, which will keep the deployer "outage" much shorter
(added by [duffj]https://github.com/duffj) on 2016-07-12).
1. Shut down neo4j's dependencies.
1. Start up neo4j's dependencies.

### Bugs

1. (MINOR) Fix the "startTime log" bug where it doesn't show the time properly.
1. (MINOR) When deployer is still active, you get two very similar-looking error messages in the log.
1. The final log message seems to be `Started Job to backup neo4j DB data files to S3.`, which comes *after* the message `Backup process complete.`, presumably because of a multi-threading ordering problem. Needs fixing.

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
[sourcegraph]: https://sourcegraph.com/github.com/Financial-Times/coco-neo4j-backup
[env-decode]: https://github.com/joeshaw/envdecode
