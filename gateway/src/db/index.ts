import Redis from 'ioredis';
import logger from '../logger';

const redis = new Redis({
    port: Number(process.env.REDIS_PORT) || 6379, // Redis port
    host: process.env.REDIS_HOST || '0.0.0.0', // Redis host
    username: process.env.REDIS_USERNAME, // needs Redis >= 6
    password: process.env.REDIS_PASSWORD,
    db: process.env.REDIS_DB ? Number(process.env.REDIS_DB) : 0, // Defaults to 0
});

const addMessage = (message: string) => {
    logger.info('hey', message);
    logger.info('hey', redis);
};

const getMessage = (id: string): string => {
    logger.info('hey', id);
    return '';
};

export default { addMessage, getMessage };
