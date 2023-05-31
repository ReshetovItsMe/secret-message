// GENERATED CODE -- DO NOT EDIT!

// package: main
// file: message.proto

import * as message_pb from "./message_pb";
import * as grpc from "@grpc/grpc-js";

interface ISecretAssistantService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  encrypt: grpc.MethodDefinition<message_pb.EncryptRequestMessage, message_pb.EncryptedMessageResponse>;
  decrypt: grpc.MethodDefinition<message_pb.DecryptRequestMessage, message_pb.DecryptedMessageResponse>;
}

export const SecretAssistantService: ISecretAssistantService;

export interface ISecretAssistantServer extends grpc.UntypedServiceImplementation {
  encrypt: grpc.handleUnaryCall<message_pb.EncryptRequestMessage, message_pb.EncryptedMessageResponse>;
  decrypt: grpc.handleUnaryCall<message_pb.DecryptRequestMessage, message_pb.DecryptedMessageResponse>;
}

export class SecretAssistantClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  encrypt(argument: message_pb.EncryptRequestMessage, callback: grpc.requestCallback<message_pb.EncryptedMessageResponse>): grpc.ClientUnaryCall;
  encrypt(argument: message_pb.EncryptRequestMessage, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<message_pb.EncryptedMessageResponse>): grpc.ClientUnaryCall;
  encrypt(argument: message_pb.EncryptRequestMessage, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<message_pb.EncryptedMessageResponse>): grpc.ClientUnaryCall;
  decrypt(argument: message_pb.DecryptRequestMessage, callback: grpc.requestCallback<message_pb.DecryptedMessageResponse>): grpc.ClientUnaryCall;
  decrypt(argument: message_pb.DecryptRequestMessage, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<message_pb.DecryptedMessageResponse>): grpc.ClientUnaryCall;
  decrypt(argument: message_pb.DecryptRequestMessage, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<message_pb.DecryptedMessageResponse>): grpc.ClientUnaryCall;
}
