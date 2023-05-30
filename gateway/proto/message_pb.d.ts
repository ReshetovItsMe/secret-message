// package: main
// file: message.proto

import * as jspb from "google-protobuf";

export class ResponseMessage extends jspb.Message {
  getBody(): Uint8Array | string;
  getBody_asU8(): Uint8Array;
  getBody_asB64(): string;
  setBody(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResponseMessage.AsObject;
  static toObject(includeInstance: boolean, msg: ResponseMessage): ResponseMessage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResponseMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResponseMessage;
  static deserializeBinaryFromReader(message: ResponseMessage, reader: jspb.BinaryReader): ResponseMessage;
}

export namespace ResponseMessage {
  export type AsObject = {
    body: Uint8Array | string,
  }
}

export class RequestMessage extends jspb.Message {
  getBody(): string;
  setBody(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestMessage.AsObject;
  static toObject(includeInstance: boolean, msg: RequestMessage): RequestMessage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RequestMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestMessage;
  static deserializeBinaryFromReader(message: RequestMessage, reader: jspb.BinaryReader): RequestMessage;
}

export namespace RequestMessage {
  export type AsObject = {
    body: string,
  }
}

