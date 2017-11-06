#!/bin/bash 
set -x
# This script spins off N drand containers and tries to verify any randomness
# produced.
# It's avery ad-hoc testing and there are probably better ways to do it but
# docker-compose had a "port being already taken" problem that I did not
# resolved...

source run_local.sh

echo "Create network $NET with subnet ${SUBNET}0/24"
docker network create "$NET" --subnet "${SUBNET}0/24"

sequence=$(seq $N -1 1)
echo "Creating and running $N nodes"
# creating the keys and compose part for each node
for i in $sequence; do
# gen key and append to group
data="$TMP/node$i/"
addr="${SUBNET}2$i:$PORT$i"
mkdir -p "$data"
echo "generating keys for $addr"
#drand keygen --keys "$data" "$addr" > /dev/null 
public="drand_id.public"
volume="$data:/root/.drand/"
allVolumes[$i]=$volume
docker run --rm --volume ${allVolumes[$i]} $IMG keygen "$addr" 
#allKeys[$i]=$data$public
cp $data$public $TMP/node$i.public
allKeys[$i]=/tmp/node$i.public
done

## generate group toml
echo $allKeys
#drand group --group "$GROUPFILE" ${allKeys[@]}
docker run --rm -v $TMP:/tmp $IMG group --group /tmp/group.toml "${allKeys[@]}"

echo "GROUP FILE:"
cat $GROUPFILE

for i in $sequence; do
# gen key and append to group
data="$TMP/node$i/"
cp $GROUPFILE "$data"drand_group.toml
cmd="run"
if [ "$i" -eq 1 ]; then
    cmd="$cmd --leader --period 2s"
fi

echo "Running docker container node$i at ${SUBNET}2$i with ${allVolumes[$i]}..."
docker run --rm --name node$i --net $NET \
            --ip ${SUBNET}2$i \
            --volume ${allVolumes[$i]} -d $IMG $cmd
sleep 0.1
done

function cleanup() {
echo "removing containers ..." 
docker rm -f $(docker ps -a -q)
}

function checkSuccess() {
    if [ "$1" -eq 0 ]; then
        return
    else
        echo "TEST <$2>: FAILURE"
        cleanup
        exit 1
    fi
}

# wait for the node to actually do the DKG and run at least one beacon
sleep 3
#docker logs node1
rootFolder="$TMP/node1"
ret=0
# check if there are any signatures
ls "$rootFolder/beacons"| grep "sig" 
checkSuccess $? "any signature produced?"

# tail returns 0 in both cases...
sigFile=$(ls "$rootFolder/beacons"| grep "sig" | tail -n 1)

# check if there is the dist public key
distPublic="$rootFolder/dist_key.public"
ls "$rootFolder/dist_key.public"
checkSuccess $? "distributed public key file?"

# try to verify with it
#drand verify --distkey "$distPublic" "$rootFolder/beacons/$sigFile"
docker run --rm -v $distPublic:/group.key -v $rootFolder/beacons/$sigFile:/beacon.sig  \
        $IMG verify --distkey /group.key  /beacon.sig
checkSuccess $? "verify signature?"

echo "TESTS OK"
cleanup
