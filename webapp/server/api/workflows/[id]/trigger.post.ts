import { defineEventHandler } from 'h3';
import { workflowQueue } from '../../../utils/queue';

export default defineEventHandler(async event => {
  const id = Number(event.context.params?.id);
  if (isNaN(id)) {
    throw createError({ statusCode: 400, statusMessage: 'Invalid workflow ID' });
  }

  const jobName = 'execute-workflow-v1';
  const jobData = { id };

  const job = await workflowQueue.add(jobName, jobData);

  return {
    message: 'Workflow trigger accepted.',
    jobId: job.id,
  };
});
