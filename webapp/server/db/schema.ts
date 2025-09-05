import { pgTable, text, serial, jsonb, timestamp, varchar, uuid } from 'drizzle-orm/pg-core';

export const workflows = pgTable('workflows', {
  id: serial('id').primaryKey(),
  name: varchar('name', { length: 256 }).notNull(),
  definition: jsonb('definition').default('{}').notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
});

export const userCredentials = pgTable('user_credentials', {
  id: uuid('id').primaryKey().defaultRandom(),
  userId: uuid('user_id').notNull(),
  service: text('service').notNull(),
  accessToken: text('access_token').notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
});
