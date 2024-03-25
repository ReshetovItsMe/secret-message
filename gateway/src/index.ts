import Hapi, { Plugin, Server } from '@hapi/hapi';
import Pino from 'hapi-pino';
import logger from './logger';
import { message } from './routes';
import Pretty from 'pino-pretty';
import HapiSwagger from 'hapi-swagger';
import Inert from '@hapi/inert';
import Vision from '@hapi/vision';

const pretty = Pretty({
    colorize: true,
});

export const init = async (): Promise<Server> => {
    const server = Hapi.server({
        port: process.env.PORT || 3000,
        host: process.env.HOST || 'localhost',
        routes: {
            cors: true,
        },
    });

    const swaggerOptions: HapiSwagger.RegisterOptions = {
        info: {
            title: 'Secret Message Gateway API Documentation',
        },
    };

    const plugins: Array<Hapi.ServerRegisterPluginObject<any>> = [
        {
            plugin: Inert,
        },
        {
            plugin: Vision,
        },
        {
            plugin: HapiSwagger,
            options: swaggerOptions,
        },
    ];

    await server.register(plugins);

    server.route(message);
    server.route({
        method: '*',
        path: '/{any*}',
        handler: (_, h) => {
            return h.response('Error! Route Not Found!').code(404);
        },
    });

    await server.register({
        plugin: Pino as Plugin<unknown>,
        options: { stream: pretty },
    });
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
