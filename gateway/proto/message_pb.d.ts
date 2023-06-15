// package: main
// file: message.proto

import * as jspb from "google-protobuf";

export class EncryptedMessageResponse extends jspb.Message {
  getEncryptedkey(): Uint8Array | string;
  getEncryptedkey_asU8(): Uint8Array;
  getEncryptedkey_asB64(): string;
  setEncryptedkey(value: Uint8Array | string): void;

  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EncryptedMessageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EncryptedMessageResponse): EncryptedMessageResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EncryptedMessageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EncryptedMessageResponse;
  static deserializeBinaryFromReader(message: EncryptedMessageResponse, reader: jspb.BinaryReader): EncryptedMessageResponse;
}

export namespace EncryptedMessageResponse {
  export type AsObject = {
    encryptedkey: Uint8Array | string,
    data: Uint8Array | string,
  }
}

export class DecryptedMessageResponse extends jspb.Message {
  getBody(): string;
  setBody(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DecryptedMessageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DecryptedMessageResponse): DecryptedMessageResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DecryptedMessageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DecryptedMessageResponse;
  static deserializeBinaryFromReader(message: DecryptedMessageResponse, reader: jspb.BinaryReader): DecryptedMessageResponse;
}

export namespace DecryptedMessageResponse {
  export type AsObject = {
    body: string,
  }
}

export class EncryptRequestMessage extends jspb.Message {
  getBody(): string;
  setBody(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EncryptRequestMessage.AsObject;
  static toObject(includeInstance: boolean, msg: EncryptRequestMessage): EncryptRequestMessage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: EncryptRequestMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EncryptRequestMessage;
  static deserializeBinaryFromReader(message: EncryptRequestMessage, reader: jspb.BinaryReader): EncryptRequestMessage;
}

export namespace EncryptRequestMessage {
  export type AsObject = {
    body: string,
  }
}

export class DecryptRequestMessage extends jspb.Message {
  getBody(): Uint8Array | string;
  getBody_asU8(): Uint8Array;
  getBody_asB64(): string;
  setBody(value: Uint8Array | string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DecryptRequestMessage.AsObject;
  static toObject(includeInstance: boolean, msg: DecryptRequestMessage): DecryptRequestMessage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DecryptRequestMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DecryptRequestMessage;
  static deserializeBinaryFromReader(message: DecryptRequestMessage, reader: jspb.BinaryReader): DecryptRequestMessage;
}

export namespace DecryptRequestMessage {
  export type AsObject = {
    body: Uint8Array | string,
  }
}

