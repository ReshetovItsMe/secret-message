#!/bin/bash

unameOut="$(uname -s)"

if [[ $unameOut == *"MINGW"* ]]
then
  PROTOC_GEN_TS_PATH="${PWD}/node_modules/.bin/protoc-gen-ts.cmd"
  PROTOC_GEN_GRPC_PATH="${PWD}/node_modules/.bin/grpc_tools_node_protoc_plugin.cmd"
else
  PROTOC_GEN_TS_PATH="${PWD}/node_modules/.bin/protoc-gen-ts"
  PROTOC_GEN_GRPC_PATH="${PWD}/node_modules/.bin/grpc_tools_node_protoc_plugin"
fi

cd ../proto || exit

OUT_DIR="../gateway/proto"
mkdir -p ${OUT_DIR}



node ../gateway/node_modules/grpc-tools/bin/protoc.js --plugin=protoc-gen-ts="${PROTOC_GEN_TS_PATH}" \
  --plugin=protoc-gen-grpc=${PROTOC_GEN_GRPC_PATH} \
  --js_out="import_style=commonjs,binary:${OUT_DIR}" \
  --ts_out="service=grpc-node,mode=grpc-js:${OUT_DIR}" \
  --grpc_out="grpc_js:${OUT_DIR}" \
  message.proto
