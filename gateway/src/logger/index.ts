import Pino from 'pino';
import Pretty from 'pino-pretty';

const stream = Pretty({
    colorize: true,
});
const logger = Pino(stream);

export default logger;
