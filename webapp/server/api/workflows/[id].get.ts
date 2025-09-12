import { defineEventHandler } from 'h3';
import { db } from '~~/server/utils/db';
import { workflows } from '~~/server/db/schema';
import { and, eq } from 'drizzle-orm';
import { serverSupabaseUser } from '#supabase/server';

export default defineEventHandler(async event => {
  const user = await serverSupabaseUser(event);
  if (!user) throw createError({ statusCode: 401, statusMessage: 'Unauthorized' });

  const workflowId = Number(event.context.params?.id);
  if (isNaN(workflowId)) throw createError({ statusCode: 400, statusMessage: 'Invalid Workflow ID' });

  const [workflow] = await db
    .select()
    .from(workflows)
    .where(and(eq(workflows.id, workflowId), eq(workflows.userId, user.id)));

  if (!workflow) throw createError({ statusCode: 404, statusMessage: 'Workflow not found' });

  return { workflow };
});
