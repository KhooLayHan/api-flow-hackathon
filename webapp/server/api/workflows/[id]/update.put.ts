import { defineEventHandler } from 'h3';
import { z } from 'zod';
import { db } from '~~/server/utils/db';
import { workflows } from '~~/server/db/schema';
import { and, eq } from 'drizzle-orm';
import { serverSupabaseUser } from '#supabase/server';

// Define a reusable schema for the workflow definition
const WorkflowDefinitionSchema = z.object({
  elements: z.array(z.any()).optional(),
});

const UpdateWorkflowSchema = z.object({
  name: z.string().min(1),
  definition: WorkflowDefinitionSchema,
  githubRepo: z.string().nullable().optional(),
});

export default defineEventHandler(async event => {
  const user = await serverSupabaseUser(event);
  if (!user) {
    throw createError({ statusCode: 401, statusMessage: 'Unauthorized' });
  }

  const workflowId = Number(event.context.params?.id);
  if (isNaN(workflowId)) {
    throw createError({ statusCode: 400, statusMessage: 'Invalid workflow ID' });
  }

  const body = await readBody(event);
  const validation = UpdateWorkflowSchema.safeParse(body);

  if (!validation.success) {
    throw createError({ statusCode: 400, statusMessage: 'Invalid request body' });
  }

  const { name, definition, githubRepo } = validation.data;

  const [updatedWorkflow] = await db
    .update(workflows)
    .set({
      name,
      definition,
      githubRepo,
      updatedAt: new Date(),
    })
    .where(and(eq(workflows.id, workflowId), eq(workflows.userId, user.id)))
    .returning();

  if (!updatedWorkflow) {
    throw createError({
      statusCode: 404,
      statusMessage: 'Workflow not found or you do not have permission to edit it.',
    });
  }

  return { workflow: updatedWorkflow };
});
