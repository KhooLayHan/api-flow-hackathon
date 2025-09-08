import { defineEventHandler, readBody, getHeader } from 'h3';
import { z } from 'zod';
import { db } from '~~/server/utils/db';
import { workflows } from '~~/server/db/migrations/schema';
import { eq } from 'drizzle-orm';
import { workflowQueue } from '~~/server/utils/queue';
import crypto from 'crypto';

// Zod schema to safely parse the parts of the GitHub payload
const GithubWebhookPayloadSchema = z.object({
  repository: z.object({
    full_name: z.string(),
  }),
  pusher: z.object({
    name: z.string(),
  }),
  head_commit: z
    .object({
      message: z.string(),
    })
    .nullable(),
});

export default defineEventHandler(async event => {
  // 1. Verify the Webhook signature for security.
  const signature = getHeader(event, 'x-hub-signature-256');
  const body = await readBody(event);

  const webhookSecret = process.env.GITHUB_WEBHOOK_SECRET;
  if (!webhookSecret) {
    console.error('GITHUB_WEBHOOK_SECRET is not set!');
    return { statusCodes: 500, message: 'Server configuration error.' };
  }

  const hmac = crypto.createHmac('sha256', webhookSecret);
  const digest = `sha256=${hmac.update(JSON.stringify(body)).digest('hex')}`;

  if (!signature || !crypto.timingSafeEqual(Buffer.from(digest), Buffer.from(signature))) {
    throw createError({ statusCode: 401, statusMessage: 'Invalid signature.' });
  }

  // 2. Parse the payload.
  const validation = GithubWebhookPayloadSchema.safeParse(body);
  if (!validation.success) {
    console.warn('Received an invalid payload or unexpected GitHub webhook payload.');
    return { statusCodes: 400, message: 'Invalid payload.' };
  }

  const payload = validation.data;
  const repoFullName = payload.repository.full_name;

  // 3. Find the matching workflow that is configured for this specific repository.
  const [workflowToTrigger] = await db.select().from(workflows).where(eq(workflows.githubRepo, repoFullName)).limit(1);
  if (!workflowToTrigger) {
    console.warn(`No workflow found for repository ${repoFullName}`);
    return { statusCodes: 200, message: 'Webhook received, no matching workflow.' };
  }

  // 4. Queue the job.
  const jobName = 'execute-workflow-v1';
  const jobData = {
    workflowId: workflowToTrigger.id,
    triggerPayload: payload,
    // userId: workflowToTrigger.userId,
  };

  await workflowQueue.add(jobName, jobData);
  console.log(`Queued job for workflow ${workflowToTrigger.id} triggered by ${repoFullName}`);

  return { statusCodes: 202, message: 'Webhook received and job queued for processing.' };
});
