import Redis from 'ioredis';
import { v4 as uuid } from 'uuid';
import logger from '../logger';

const redis = new Redis({
    port: Number(process.env.REDIS_PORT) || 6379, // Redis port
    host: process.env.REDIS_HOST || '0.0.0.0', // Redis host
    username: process.env.REDIS_USERNAME, // needs Redis >= 6
    password: process.env.REDIS_PASSWORD,
    db: process.env.REDIS_DB ? Number(process.env.REDIS_DB) : 0, // Defaults to 0
});

const addMessage = (message: string): string => {
    logger.info('Adding message');
    const id = uuid();
    logger.trace(`Message id ${id}`);
    redis.set(uuid(), message);
    return id;
};

const getMessage = async (id: string): Promise<string> => {
    logger.trace(`Get message id ${id}`);
    const message = await redis.get(id);
    if (message) {
        return message;
    }
    throw new Error(`Message with id=${id} doesn't exist`);
};

export default { addMessage, getMessage };
