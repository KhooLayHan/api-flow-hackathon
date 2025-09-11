import { defineEventHandler, readBody } from 'h3';
import { z } from 'zod';
import { db } from '~~/server/utils/db';
import { workflows } from '~~/server/db/schema';
import { serverSupabaseUser } from '#supabase/server';

const SaveWorkflowSchema = z.object({
  name: z.string().min(1, { message: 'Workflow name cannot be empty' }),
  definition: z.record(z.any(), z.unknown()),
});

export default defineEventHandler(async event => {
  const user = await serverSupabaseUser(event);
  if (!user) {
    throw createError({
      statusCode: 401,
      statusMessage: 'Unauthorized',
    });
  }

  const body = await readBody(event);
  const validation = SaveWorkflowSchema.safeParse(body);

  if (!validation.success) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Invalid request body',
      data: validation.error.issues,
    });
  }

  const { name, definition } = validation.data;
  console.log(`Saving workflow: ${name}`, definition);

  const [newWorkflow] = await db.insert(workflows).values({ name, definition }).returning();

  return {
    workflow: newWorkflow,
  };
});
