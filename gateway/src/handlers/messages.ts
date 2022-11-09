import { Lifecycle } from '@hapi/hapi';
import logger from '../logger';
import db from '../db';

const sendMessage: Lifecycle.Method = (request, h) => {
    try {
        logger.info('Got new request for add message');
        const { payload } = request;
        const { message } = payload as { message: string };

        logger.info('Start encryption');
        const encryptedMessage = message;
        logger.info('Encrypted');

        const id = db.addMessage(encryptedMessage);

        return h.response({ messageId: id }).code(201);
    } catch (error: unknown) {
        return h.response(error as Error).code(500);
    }
};

export default { sendMessage };
