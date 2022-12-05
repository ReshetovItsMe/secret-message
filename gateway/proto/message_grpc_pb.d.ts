// GENERATED CODE -- DO NOT EDIT!

// package: main
// file: message.proto

import * as message_pb from "./message_pb";
import * as grpc from "@grpc/grpc-js";

interface ISecretAssistantService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  encrypt: grpc.MethodDefinition<message_pb.RequestMessage, message_pb.ResponseMessage>;
}

export const SecretAssistantService: ISecretAssistantService;

export interface ISecretAssistantServer extends grpc.UntypedServiceImplementation {
  encrypt: grpc.handleUnaryCall<message_pb.RequestMessage, message_pb.ResponseMessage>;
}

export class SecretAssistantClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  encrypt(argument: message_pb.RequestMessage, callback: grpc.requestCallback<message_pb.ResponseMessage>): grpc.ClientUnaryCall;
  encrypt(argument: message_pb.RequestMessage, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<message_pb.ResponseMessage>): grpc.ClientUnaryCall;
  encrypt(argument: message_pb.RequestMessage, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<message_pb.ResponseMessage>): grpc.ClientUnaryCall;
}
