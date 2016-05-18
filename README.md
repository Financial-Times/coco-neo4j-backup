coco-neo4j-backup
=================

Docker Image for automated neo4j backups to S3.


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
