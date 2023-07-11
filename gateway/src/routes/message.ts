import { ServerRoute } from '@hapi/hapi';
import Joi from 'joi';
import { messageHandler } from '../handlers';

const messageRoutes: Array<ServerRoute> = [
    {
        method: 'POST',
        path: '/message',
        handler: messageHandler.sendMessage,
        options: {
            description: 'Encrypt message',
            notes: 'Return a message id',
            tags: ['api'],
            validate: {
                payload: Joi.object(),
            },
        },
    },
    {
        method: 'GET',
        path: '/message',
        handler: messageHandler.getMessage,
        options: {
            description: 'Decrypt message',
            notes: 'Return decrypted message',
            tags: ['api'],
            validate: {
                query: Joi.object({
                    messageId: Joi.string(),
                }),
            },
        },
    },
];

export default messageRoutes;
