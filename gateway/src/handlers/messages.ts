import { Lifecycle } from '@hapi/hapi';
import logger from '../logger';
import { db, secretAssistant } from '../services';
import { RequestMessage } from '../../proto/message_pb';

const sendMessage: Lifecycle.Method = async (request, h) => {
    try {
        logger.info('Got new request for add message');
        const { payload } = request;
        const { message } = payload as { message: string };

        logger.info('Start encryption');
        const reqMessage: RequestMessage = new RequestMessage();
        reqMessage.setBody(message);
        const encryptedMessage: string = await new Promise(
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
                    resolve(newMessage.getBody());
                });
            }
        );
        logger.info('Encrypted');

        const id = db.addMessage(encryptedMessage);

        return h.response({ messageId: id }).code(201);
    } catch (error: unknown) {
        const err = error as Error;
        logger.error('Error, something went wrong:');
        logger.error(err);
        return h.response(err.message).code(500);
    }
};

export default { sendMessage };
