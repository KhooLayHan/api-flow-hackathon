import { defineEventHandler } from 'h3';
import { db } from '../../utils/db';
import { workflows } from '../../db/schema';
import { eq } from 'drizzle-orm';

export default defineEventHandler(async event => {
  const id = Number(event.context.params?.id);
  if (isNaN(id)) throw createError({ statusCode: 400, statusMessage: 'Invalid ID' });

  const [workflow] = await db.select().from(workflows).where(eq(workflows.id, id));
  if (!workflow) throw createError({ statusCode: 404, statusMessage: 'Workflow not found' });

  return { workflow };
});
