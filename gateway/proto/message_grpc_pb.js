// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var message_pb = require('./message_pb.js');

function serialize_main_DecryptRequestMessage(arg) {
  if (!(arg instanceof message_pb.DecryptRequestMessage)) {
    throw new Error('Expected argument of type main.DecryptRequestMessage');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_main_DecryptRequestMessage(buffer_arg) {
  return message_pb.DecryptRequestMessage.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_main_DecryptedMessageResponse(arg) {
  if (!(arg instanceof message_pb.DecryptedMessageResponse)) {
    throw new Error('Expected argument of type main.DecryptedMessageResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_main_DecryptedMessageResponse(buffer_arg) {
  return message_pb.DecryptedMessageResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_main_EncryptRequestMessage(arg) {
  if (!(arg instanceof message_pb.EncryptRequestMessage)) {
    throw new Error('Expected argument of type main.EncryptRequestMessage');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_main_EncryptRequestMessage(buffer_arg) {
  return message_pb.EncryptRequestMessage.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_main_EncryptedMessageResponse(arg) {
  if (!(arg instanceof message_pb.EncryptedMessageResponse)) {
    throw new Error('Expected argument of type main.EncryptedMessageResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_main_EncryptedMessageResponse(buffer_arg) {
  return message_pb.EncryptedMessageResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var SecretAssistantService = exports.SecretAssistantService = {
  encrypt: {
    path: '/main.SecretAssistant/encrypt',
    requestStream: false,
    responseStream: false,
    requestType: message_pb.EncryptRequestMessage,
    responseType: message_pb.EncryptedMessageResponse,
    requestSerialize: serialize_main_EncryptRequestMessage,
    requestDeserialize: deserialize_main_EncryptRequestMessage,
    responseSerialize: serialize_main_EncryptedMessageResponse,
    responseDeserialize: deserialize_main_EncryptedMessageResponse,
  },
  decrypt: {
    path: '/main.SecretAssistant/decrypt',
    requestStream: false,
    responseStream: false,
    requestType: message_pb.DecryptRequestMessage,
    responseType: message_pb.DecryptedMessageResponse,
    requestSerialize: serialize_main_DecryptRequestMessage,
    requestDeserialize: deserialize_main_DecryptRequestMessage,
    responseSerialize: serialize_main_DecryptedMessageResponse,
    responseDeserialize: deserialize_main_DecryptedMessageResponse,
  },
};

exports.SecretAssistantClient = grpc.makeGenericClientConstructor(SecretAssistantService);
