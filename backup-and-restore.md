Backup and Restore
==================

Running a Backup
----------------

### Things to note

* At the time of writing, the neo4j backup process is not automated or scheduled. Making it so is in the 
* The backup process is tied to the red neo4j instance. Making it configurable to red or blue is in the backlog - see below.
* At the time of writing, the overall process takes around 40min.
* The backlog is at [README](README.md).

### Requirements

To run a backup (using the `semantic` cluster as an example), you will need the following:

* A basic understanding of Linux, the console and how to execute SSH commands, with the appropriate software (PuTTY, Linux/Mac) to do so.
* Permission to access the cluster whose database you are trying to back up. (Trying to execute the first SSH command should tell you whether you have this or not.)
* Ideally, access to the AWS InfraProd account, so that you can verify that your backup has been created.
* Ideally, a basic understanding of CoCo and the UPP stack, so that you are aware of the rough impact of the commands that you are executing.

### Steps

1. SSH to the cluster:

        ssh core@semantic-tunnel-up.ft.com

1. Warn people in the appropriate Slack channel (e.g. [#co-co](https://financialtimes.slack.com/messages/co-co/)) that you are about to stop the deployer in that cluster and run a backup.

1. Once you are satisfied that people are happy that you are running a backup, stop the deployer, red ingester, and red annotators, then run the backup and watch the logs (it should take about half an hour):

        fleetctl stop deployer.service content-ingester-neo4j-red@1.service v1-content-annotator-red@1.service v2-content-annotator-red@1.service \
            && fleetctl start neo4j-backup.service \
            && fleetctl journal -f neo4j-backup.service

1. Start the services that you stopped earlier:

        fleetctl start deployer.service content-ingester-neo4j-red@1.service v1-content-annotator-red@1.service v2-content-annotator-red@1.service

1. *(OPTIONAL)* Verify the backup by downloading it to your local machine, extracting it and starting up a local Neo4j instance pointing to the backed up data, then playing with the data until you are satisfied that it is complete.

### Tips

* While the rsync process is running, a directory will be growing inside `/vol/neo4j-red-1` on the host machine running `neo4j-backup.service`. To monitor it, do this:

        fleetctl ssh neo4j-backup.service
        watch du -hs /vol/neo4j-red-1
        
* Once the backup has started streaming up to S3, you can see it by logging in to the [InfraProd](https://awslogin.internal.ft.com/InfraProd/default.aspx) FT AWS account and looking at the [com.ft.universalpublishing.backup-data](https://console.aws.amazon.com/s3/home?region=eu-west-1#&bucket=com.ft.universalpublishing.backup-data&prefix=) bucket. The start of the process is accompanied by this sort of log message in journald:

        Jul 12 08:49:34 ip-172-24-151-154.eu-west-1.compute.internal systemd[1]: Started Job to backup neo4j DB data files to S3.

* You might want to keep an eye on the output of `fleetctl list-machines`, because in the past, a backup has caused one of the host machines to crash; this problem should have been resolved by the addition of `nice` to the rsync process. If `nice` isn't enough, there is a feature in the TODO list below to also add `ionice`.


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

1. Exit the container and the VM.

1. The following instructions need to be repeated for the `red` and `blue` neo4j instances. To ease the process, an environment variable
has been used which contains the 'colour' of the neo4j instance.

    1. Set the `$NEO_COLOUR` environment variable:

            export NEO_COLOUR=red # then repeat with 'blue'

    1. SSH to the aforementioned red cluster host and re-export the envvar:

            fleetctl ssh neo4j-${NEO_COLOUR}@1.service
            export NEO_COLOUR=red # then repeat with 'blue'
            export ARCHIVE_NAME="/vol/neo4j_backup_2016-06-06T12-55-09_semantic.tar.gz"

    1. Stop dependent services:
    
            fleetctl stop content-ingester-neo4j-${NEO_COLOUR}@1.service v1-content-annotator-${NEO_COLOUR}@1.service v2-content-annotator-${NEO_COLOUR}@1.service

    1. Stop neo:
    
            fleetctl stop neo4j-${NEO_COLOUR}@1.service

    1. SCP the backup file to the cluster host that was running the red neo4j instance (SCP only needed if you aren't where you downloaded the tarball to)
    1. Back up the old data directory and extract the contents of the backup tarball into the `/vol` partition (please note that the expected tar ball is from the backup service and may need to be adjusted for any other way in which the tar file may have been created):
    
            cd /vol/neo4j-${NEO_COLOUR}-1 \
                && ls -l \
                && sudo mv graph.db graph.db.old.`date +%s` \
                && sudo mv graph.db.backup graph.db.backup.old.`date +%s` \
                && sudo tar -xzvf $ARCHIVE_NAME --strip-components=1 \
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

