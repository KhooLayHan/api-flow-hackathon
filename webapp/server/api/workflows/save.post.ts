import { defineEventHandler, readBody } from 'h3';
import { z } from 'zod';

const SaveWorkflowSchema = z.object({
  name: z.string().min(1, { message: 'Workflow name cannot be empty' }),

  definition: z.record(z.never(), z.unknown()),

  workflowId: z.uuid().optional(),
  // description: z.string().min(2).max(500).optional(),
  // steps: z.array(
  //   z.object({
  //     id: z.string().uuid(),
  //     name: z.string().min(2).max(100),
  //     description: z.string().min(2).max(500).optional(),
  //     type: z.string().min(2).max(100),
  //     config: z.any().optional(),
  //   })
  // ),
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

  const { name, definition, workflowId } = validation.data;

  console.log(`Saving workflow: ${name}`, definition);

  return {
    status: 'ok',
    message: 'Workflow saved successfully',
  };
  // Save the workflow to the database
  // const workflow = await saveWorkflow(name, definition, workflowId);

  // return workflow;
});
