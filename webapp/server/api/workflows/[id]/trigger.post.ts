import { defineEventHandler } from 'h3';
import { workflowQueue } from '../../../utils/queue';
import { serverSupabaseUser } from '#supabase/server';

export default defineEventHandler(async event => {
  const user = await serverSupabaseUser(event);
  if (!user) {
    throw createError({ statusCode: 401, statusMessage: 'Unauthorized' });
  }

  const workflowId = Number(event.context.params?.id);
  if (isNaN(workflowId)) {
    throw createError({ statusCode: 400, statusMessage: 'Invalid workflow ID' });
  }

  const jobName = 'execute-workflow-v1';
  // Payload for the Go worker
  const jobData = { workflowId: workflowId, userId: user.id };

  const job = await workflowQueue.add(jobName, jobData);

  return {
    message: 'Workflow trigger accepted.',
    jobId: job.id,
  };
});
