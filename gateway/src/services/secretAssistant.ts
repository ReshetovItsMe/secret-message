import { credentials } from '@grpc/grpc-js';
import { SecretAssistantClient } from '../../proto/message_grpc_pb';

const client = new SecretAssistantClient(
    process.env.SECRET_ASSISTANT_URL ?? 'secret-assistant:50052',
    credentials.createInsecure()
);

export default client;
