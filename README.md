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

1. Stop the deployer, red ingester, and red annotators, then run the backup and watch the logs (it should take about half an hour):

        fleetctl stop deployer.service content-ingester-neo4j-red@1.service v1-content-annotator-red@1.service v2-content-annotator-red@1.service \
            && fleetctl start neo4j-backup.service \
            && fleetctl journal -f neo4j-backup.service

1. Start the services that you stopped earlier:

        fleetctl start deployer.service content-ingester-neo4j-red@1.service v1-content-annotator-red@1.service v2-content-annotator-red@1.service


How to run a restore
--------------------

This is a cold process, so neo4j and its services will be unavailable while this is happening.

1. Warn everybody.
1. SSH into the cluster:

        ssh core@semantic-tunnel-up.ft.com

1. Stop the deployer:

        fleetctl stop deployer.service

1. Download the backup file from S3 (~4GB at the time of writing) by firing up a container that has the AWS command-line client installed:

    1. SSH to the machine running neo and create/enter a container with the `/vol` volume mounted:

            fleetctl ssh neo4j-red@1.service
            docker run -ti \
                --env "AWS_ACCESS_KEY_ID=$(etcdctl get /ft/_credentials/aws/aws_access_key_id)" \
                --env "AWS_SECRET_ACCESS_KEY=$(etcdctl get /ft/_credentials/aws/aws_secret_access_key)" \
                --env "AWS_DEFAULT_REGION=eu-west-1" \
                -v /vol:/vol alpine sh

    1. Inside the container, export AWS credentials as environment variables, then install AWS command line tools:

            apk --update add python py-pip
            pip install awscli
            
    1. Finally, copy the tarball from S3 to `/vol`:

            cd /vol
            aws s3 cp s3://com.ft.universalpublishing.backup-data/neo4j_backup_2016-06-06T12-55-09_semantic.tar.gz .

1. Exit the container and the VM

1. The following instructions need to be repeated for the `red` and `blue` neo4j instances. To ease the process, an environment variable
has been used which contains the 'colour' of the neo4j instance.

    1. Set the `$NEO_COLOUR` environment variable:

            export NEO_COLOUR=red # then repeat with 'blue'

    1. SSH to the aforementioned red cluster host and re-export the envvar:

            fleetctl ssh neo4j-${NEO_COLOUR}@1.service
            export NEO_COLOUR=red # then repeat with 'blue'
            export ARCHIVE_NAME="neo4j_backup_2016-06-06T12-55-09_semantic.tar.gz"

    1. Stop dependent services:
    
            fleetctl stop content-ingester-neo4j-${NEO_COLOUR}@1.service v1-content-annotator-${NEO_COLOUR}@1.service v2-content-annotator-${NEO_COLOUR}@1.service

    1. Stop neo:
    
            fleetctl stop neo4j-${NEO_COLOUR}@1.service

    1. SCP the backup file to the cluster host that was running the red neo4j instance (SCP only needed if you aren't where you downloaded the tarball to)
    1. Back up the old data directory and extract the contents of the backup tarball into the `/vol` partition (please note that the expected tar ball is from the backup service and may need to be adjusted for any other way in which the tar file may have been created):
    
            cd /vol/neo4j-${NEO_COLOUR}-1 \
                && sudo mv graph.db graph.db.old \
                && sudo mv graph.db.backup graph.db.backup.old \
                && sudo tar -xzvf /vol/neo4j-${NEO_COLOUR}-1/$ARCHIVE_NAME --strip-components=1 \
                && sudo mv graph.db.backup graph.db

    1. Start up neo and its dependencies:

            fleetctl start neo4j-${NEO_COLOUR}@1.service \
                && sleep 10 \
                && fleetctl start content-ingester-neo4j-${NEO_COLOUR}@1.service v1-content-annotator-${NEO_COLOUR}@1.service v2-content-annotator-${NEO_COLOUR}@1.service

    1. Test that everything is ok, e.g. by getting some sample data from the red API (you will need to flesh the below curl command
    out with the rest of the parameters required:

            curl -H "<auth header>" https://semantic-up.ft.com/__people-rw-neo4j-${NEO_COLOUR}/people/{uuid}

1. Repeat the process for the blue system, starting with:

        exit
        fleetctl ssh neo4j-blue@1.service
        export NEO_COLOUR=blue

1. Once you're happy that the restore process has worked properly, you can start the deployer and let everyone know that the world
is ok again:

        fleetctl start deployer.service


NB. It may be possible to condense the S3 download process down into a single command, but at the time of writing it failed to stream
the archive file properly from S3 so I have left it out of the official instructions:
 
    docker run -ti  -v '/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt' -v '/vol:/vol' \
        --env "AWS_ACCESS_KEY_ID=$(etcdctl get /ft/_credentials/aws/aws_access_key_id)" \
        --env "AWS_SECRET_ACCESS_KEY=$(etcdctl get /ft/_credentials/aws/aws_secret_access_key)" \
        coco/gof3r gof3r get -b com.ft.universalpublishing.backup-data -k neo4j_backup_2016-06-06T12-55-09_semantic.tar.gz \
            --endpoint s3-eu-west-1.amazonaws.com > /vol/tmp/neo4j_backup_2016-06-06T12-55-09_semantic.tar.gz


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

1. Shut down neo4j's dependencies.
1. Start up neo4j's dependencies.
1. Shamelessly plagiarise mongo-backup.timer to create neo4j-backup.timer.
1. Stop and start the deployer programmatically to avoid neo4j being accidentally started up during a backup.
1. Upload backups into a folder inside the bucket.
1. Write a health check.
1. ~~Lock down the version in services.yaml to a specific tag.~~ DONE
1. Write more tests. Always more tests.
1. Print a link to the backup archive in S3.
1. Check CPU usage, then see if using an LZ4 compressor reduces CPU usage (potentially at the cost of a larger backup file).
1. Switch to using a library like [env-decode] for much simpler parsing of environment variables without needing CLI params, which
are unnecessary for most apps.

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
