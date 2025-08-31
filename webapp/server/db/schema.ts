import { pgTable, serial, jsonb, timestamp, varchar } from 'drizzle-orm/pg-core';

export const workflows = pgTable('workflows', {
  id: serial('id').primaryKey(),
  name: varchar('name', { length: 256 }).notNull(),
  definition: jsonb('definition').default('{}').notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
});
