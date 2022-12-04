// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var message_pb = require('./message_pb.js');

function serialize_main_RequestMessage(arg) {
  if (!(arg instanceof message_pb.RequestMessage)) {
    throw new Error('Expected argument of type main.RequestMessage');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_main_RequestMessage(buffer_arg) {
  return message_pb.RequestMessage.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_main_ResponseMessage(arg) {
  if (!(arg instanceof message_pb.ResponseMessage)) {
    throw new Error('Expected argument of type main.ResponseMessage');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_main_ResponseMessage(buffer_arg) {
  return message_pb.ResponseMessage.deserializeBinary(new Uint8Array(buffer_arg));
}


var SecretAssistantService = exports.SecretAssistantService = {
  encrypt: {
    path: '/main.SecretAssistant/encrypt',
    requestStream: false,
    responseStream: false,
    requestType: message_pb.RequestMessage,
    responseType: message_pb.ResponseMessage,
    requestSerialize: serialize_main_RequestMessage,
    requestDeserialize: deserialize_main_RequestMessage,
    responseSerialize: serialize_main_ResponseMessage,
    responseDeserialize: deserialize_main_ResponseMessage,
  },
};

exports.SecretAssistantClient = grpc.makeGenericClientConstructor(SecretAssistantService);
