#!/bin/bash
set -mx
ROOT_PROJECT_NAME="hcc"
PROJECT_NAME="cello"

PROTO_PROJECT_NAME="melody"
PACKAGING_SCRIPT_FILE="packaging.sh"

if [ $(uname -s) == "FreeBSD" ]; then
 exit 0
fi



#  cp -r `echo $GOPATH`/src/$ROOT_PROJECT_NAME/$PROTO_PROJECT_NAME ./tmp_$PROTO_PROJECT_NAME
cp -r $GOPATH/src/$ROOT_PROJECT_NAME/$PROTO_PROJECT_NAME ./tmp_$PROTO_PROJECT_NAME
./tmp_$PROTO_PROJECT_NAME/$PACKAGING_SCRIPT_FILE $PROJECT_NAME


protoc -I /usr/local/include -I  ./tmp_$PROTO_PROJECT_NAME --go_out=$GOPATH/src --go-grpc_out=$GOPATH/src ./tmp_$PROTO_PROJECT_NAME/*.proto
rm -rf ./tmp_$PROTO_PROJECT_NAME

