#!/bin/bash

#PROTO_REPO="https://github.com/your-org/grpc-protos"
PROTO_DIR="../rpc/protos"
#TEMP_DIR="./temp_protos"
#
## Remove existing proto files
#rm -rf $PROTO_DIR
#mkdir -p $PROTO_DIR
#
## Clone the proto repository
#git clone $PROTO_REPO $TEMP_DIR
#
## Copy proto files to the desired directory
#cp $TEMP_DIR/*.proto $PROTO_DIR
#
## Remove the temporary directory
#rm -rf $TEMP_DIR

# Generate gRPC code
protoc --go_out=$PROTO_DIR --go-grpc_out=$PROTO_DIR --proto_path=$PROTO_DIR $PROTO_DIR/*.proto
