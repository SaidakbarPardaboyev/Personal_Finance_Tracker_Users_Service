#!/bin/bash
CURRENT_DIR=$1
rm -rf ${CURRENT_DIR}/genproto
for x in $(find ${CURRENT_DIR}/Personal_Finance_Tracker_Protos/* -type d); do
  protoc -I=${x} -I=${CURRENT_DIR}/Personal_Finance_Tracker_Protos -I /usr/local/go --go_out=${CURRENT_DIR} \
   --go-grpc_out=${CURRENT_DIR} ${x}/*.proto
done
