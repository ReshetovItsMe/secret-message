import Redis from 'ioredis';
import { v4 as uuid } from 'uuid';
import logger from '../logger';

const redis = new Redis({
    port: Number(process.env.REDIS_PORT) || 6379, // Redis port
    host: process.env.REDIS_HOST || 'redis', // Redis host
    username: process.env.REDIS_USERNAME, // needs Redis >= 6
    password: process.env.REDIS_PASSWORD,
    db: process.env.REDIS_DB ? Number(process.env.REDIS_DB) : 0, // Defaults to 0
});

const addEncryptionKey = async (key: Uint8Array): Promise<string> => {
    logger.info('Adding AES encryption key');
    const id = uuid();
    await redis.set(id, Buffer.from(key));
    return id;
};

const addData = async (id: Uint8Array, key: Uint8Array): Promise<void> => {
    logger.info('Adding private encryption key');
    await redis.set(Buffer.from(id), Buffer.from(key));
};

const getData = async (id: string | Buffer): Promise<Buffer> => {
    logger.trace(`Get message id ${id}`);
    const message = await redis.getBuffer(id);
    if (message) {
        return message;
    }
    throw new Error(`Message with id=${id} doesn't exist`);
};

export default { addData, getData, addEncryptionKey };
