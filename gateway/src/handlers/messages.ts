import { Lifecycle } from '@hapi/hapi';
import logger from '../logger';
import { db, secretAssistant } from '../services';
import {
    EncryptRequestMessage,
    DecryptRequestMessage,
} from '../../proto/message_pb';

const sendMessage: Lifecycle.Method = async (request, h) => {
    try {
        logger.info('Got new request for add message');
        const { payload } = request;
        const { message } = payload as { message: string };

        logger.info('Start encryption');
        const reqMessage: EncryptRequestMessage = new EncryptRequestMessage();
        reqMessage.setBody(message);

        const encryptedMessage: Uint8Array = await new Promise(
            (resolve, reject) => {
                secretAssistant.encrypt(reqMessage, (error, newMessage) => {
                    if (error) {
                        reject(error);
                        return;
                    }
                    if (!newMessage) {
                        reject(new Error('Encrypted message is empty'));
                        return;
                    }
                    resolve(newMessage.getBody() as Uint8Array);
                });
            }
        );
        logger.info('Encrypted');

        const id = await db.addMessage(encryptedMessage);

        return h.response({ messageId: id }).code(201);
    } catch (error: unknown) {
        const err = error as Error;
        logger.error('Error, something went wrong:');
        logger.error(err);
        return h.response(err.message).code(500);
    }
};

const getMessage: Lifecycle.Method = async (request, h) => {
    try {
        logger.info('Got new request for get message');
        const { query } = request;
        const { messageId } = query as { messageId: string };

        logger.info('Start decryption');

        const message = await db.getMessage(messageId);

        const reqMessage: DecryptRequestMessage = new DecryptRequestMessage();
        reqMessage.setBody(message);

        const decryptedMessage: Uint8Array | string = await new Promise(
            (resolve, reject) => {
                secretAssistant.decrypt(reqMessage, (error, newMessage) => {
                    if (error) {
                        reject(error);
                        return;
                    }
                    if (!newMessage) {
                        reject(new Error('Encrypted message is empty'));
                        return;
                    }
                    resolve(newMessage.getBody());
                });
            }
        );
        logger.info('Decrypted');

        return h.response({ message: decryptedMessage }).code(201);
    } catch (error: unknown) {
        const err = error as Error;
        logger.error('Error, something went wrong:');
        logger.error(err);
        return h.response(err.message).code(500);
    }
};

export default { sendMessage, getMessage };
