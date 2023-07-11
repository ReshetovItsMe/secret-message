import { Lifecycle } from '@hapi/hapi';
import logger from '../logger';
import { db, secretAssistant } from '../services';
import {
    EncryptRequestMessage,
    DecryptRequestMessage,
    EncryptedMessageResponse,
} from '../../proto/message_pb';

interface DecodedEncryptedMessage {
    privateKey: Uint8Array;
    encryptedKey: Uint8Array;
    data: Uint8Array;
}

const sendMessage: Lifecycle.Method = async (request, h) => {
    try {
        logger.info('Got new request for add message');
        const { payload } = request;
        const { message } = payload as { message: string };

        logger.info('Start encryption');
        const reqMessage: EncryptRequestMessage = new EncryptRequestMessage();
        reqMessage.setBody(message);

        const encryptedMessage: EncryptedMessageResponse = await new Promise(
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
                    resolve(newMessage);
                });
            }
        );
        logger.info('Encrypted');

        const jsonString = Buffer.from(encryptedMessage.getBody()).toString();

        let decodedMessage: DecodedEncryptedMessage = JSON.parse(jsonString);

        decodedMessage = {
            encryptedKey: new Uint8Array(
                Buffer.from(decodedMessage.encryptedKey)
            ),
            privateKey: new Uint8Array(Buffer.from(decodedMessage.privateKey)),
            data: new Uint8Array(Buffer.from(decodedMessage.data)),
        };

        const id = await db.addEncryptionKey(decodedMessage.encryptedKey);
        await db.addData(
            decodedMessage.encryptedKey,
            decodedMessage.privateKey
        );
        await db.addData(decodedMessage.privateKey, decodedMessage.data);

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

        const encryptedKey = await db.getData(messageId);
        const privateKey = await db.getData(encryptedKey);
        const message = await db.getData(privateKey);

        const reqMessage: DecryptRequestMessage = new DecryptRequestMessage();

        const body = {
            privateKey: privateKey.toString(),
            encryptedKey: encryptedKey.toString(),
            data: message.toString(),
        };

        reqMessage.setBody(Buffer.from(JSON.stringify(body)));

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
