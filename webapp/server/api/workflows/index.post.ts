import { defineEventHandler, readBody } from 'h3';
import { z } from 'zod';
import { db } from '../../utils/db';
import { workflows } from '../../db/schema';

const SaveWorkflowSchema = z.object({
  name: z.string().min(1, { message: 'Workflow name cannot be empty' }),
  definition: z.record(z.never(), z.unknown()),
});

export default defineEventHandler(async event => {
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
