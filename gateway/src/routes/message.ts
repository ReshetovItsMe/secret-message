import { ServerRoute } from '@hapi/hapi';
import { messageHandler } from '../handlers';

const messageRoutes: Array<ServerRoute> = [
    {
        method: 'POST',
        path: '/messages',
        handler: messageHandler.sendMessage,
    },
];

export default messageRoutes;
