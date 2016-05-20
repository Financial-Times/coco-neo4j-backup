coco-neo4j-backup
=================

Docker Image for automated neo4j backups to S3.

**NB. As of 2016-05-20, this is a work in progress.**

Requirements
------------

* runs in CoCo
* runs in a different container, but mounts the neo container's filesystem
so that it can do rsync and all that jazz
* uses fleet Go API to shut down neo4j
* needs to tell the deployer not to auto-restart neo4j; there is currently
no API for this.


TODO
----

1. Set up a skeleton codebase.
1. Write the code to:

    1. Shut down neo4j's dependencies.
    1. Shut down neo4j.
    1. Create a backup artefact using tar and gzip.
    1. Upload the archive to S3.
    1. Start up neo4j.
    1. Start up neo4j's dependencies.


Notes and Questions
-------------------

1. This thing has to run on the same box as neo4j, right? Is that possible/easy to do in a container-based world?

    * A: Actually, it'll run in its own container and mount the neo4j volume to access the files.
    
2. Why does my IntelliJ build fail when I try to access a function in another file in the same directory? See below:

        $ go build -o "/private/var/folders/rt/8c3952t54cd5q7x08z4m6j5m0000gn/T/Build main.go and rungo" /Users/dafydd/dev/go/src/github.com/Financial-Times/coco-neo4j-backup/main.go
        src/github.com/Financial-Times/coco-neo4j-backup/main.go:24: undefined: backup

    * A: 