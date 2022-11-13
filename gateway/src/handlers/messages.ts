import { Lifecycle } from '@hapi/hapi';
import logger from '../logger';
import { db, secretAssistant } from '../services';

const sendMessage: Lifecycle.Method = async (request, h) => {
    try {
        logger.info('Got new request for add message');
        const { payload } = request;
        const { message } = payload as { message: string };

        logger.info('Start encryption');
        const encryptedMessage: string = await new Promise(
            (resolve, reject) => {
                secretAssistant.encrypt(
                    Object.create({ body: message }),
                    (error, newMessage) => {
                        if (error) {
                            reject(error);
                        }
                        if (!newMessage) {
                            throw new Error('Encrypted message is empty');
                        }
                        resolve(newMessage.getBody());
                    }
                );
            }
        );
        logger.info('Encrypted');

        const id = db.addMessage(encryptedMessage);

        return h.response({ messageId: id }).code(201);
    } catch (error: unknown) {
        return h.response(error as Error).code(500);
    }
};

export default { sendMessage };
