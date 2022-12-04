import grpc from '@grpc/grpc-js';
import { SecretAssistantClient } from '../../../proto/generated/ts/message_grpc_pb';

const client = new SecretAssistantClient(
    process.env.SECRET_ASSISTANT_URL ?? 'localhost:50051',
    grpc.credentials.createInsecure()
);

export default client;
