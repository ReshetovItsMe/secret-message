import Hapi, { Server } from '@hapi/hapi';
import Pino from 'hapi-pino';
import logger from './logger';
import { message } from './routes';

export const init = async (): Promise<Server> => {
    const server = Hapi.server({
        port: process.env.PORT || 3000,
        host: process.env.HOST || 'localhost',
    });

    server.route(message);
    server.route({
        method: '*',
        path: '/{any*}',
        handler: (_, h) => {
            return h.response('Error! Route Not Found!').code(404);
        },
    });

    await server.register(Pino);
    server.log(`Listening on ${server.settings.host}:${server.settings.port}`);

    await server.start();
    return server;
};

process.on('unhandledRejection', (err) => {
    logger.error('unhandledRejection');
    logger.error(err);
    process.exit(1);
});

init();
