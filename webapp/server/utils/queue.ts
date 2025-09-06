import { Queue } from 'bullmq';
import IORedis from 'ioredis';

const redisURL = new URL(process.env.REDIS_URL!);

const connection = new IORedis(process.env.REDIS_URL!, {
  maxRetriesPerRequest: null,
});

export const workflowQueue = new Queue('workflows', {
  connection,
});
