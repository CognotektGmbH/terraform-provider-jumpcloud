#!/bin/bash
if [ -z "$1" ]
    then
    echo "Please provide the release version as the 1st command line arg"
    exit 1
fi

RELEASE_DIR="/tmp/terraform-provider-jumpcloud/releases/$1"

echo "Creating source archive"
tar -cvf "${RELEASE_DIR}/source-{$1}.tar.gz" ./

mkdir -p "${RELEASE_DIR}"
for GOARCH in 386 amD64
do
  for GOOS in darwin linux
  do
    echo "Creating release binary for ${GOOS} / ${GOARCH}"
  	env GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-jumpcloud
    TARGET_FILENAME="${RELEASE_DIR}/terraform-provider-jumpcloud_${1}_${GOOS}_${GOARCH}.tar.gz"
    tar -cvf "${TARGET_FILENAME}" terraform-provider-jumpcloud
    TARGETS="$TARGETS $TARGET_FILENAME"
  done
done

echo "Creating md5 checksum"
md5sum "${TARGETS}" > "${RELEASE_DIR}/checksum-md5.txt"

echo "Creating sh1 checksum"
sha1sum "${TARGETS}" > "${RELEASE_DIR}/checksum-sha1.txt"
