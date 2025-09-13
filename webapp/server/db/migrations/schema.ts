import { pgTable, serial, varchar, jsonb, timestamp, uuid, text } from "drizzle-orm/pg-core"
// import { sql } from "drizzle-orm"

export const workflows = pgTable("workflows", {
  id: serial().primaryKey().notNull(),
	name: varchar({ length: 256 }).notNull(),
	definition: jsonb().default({}).notNull(),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow().notNull(),
	updatedAt: timestamp("updated_at", { mode: 'string' }).defaultNow().notNull(),
	githubRepo: varchar("github_repo", { length: 256 }).notNull(),
});

export const userCredentials = pgTable("user_credentials", {
	id: uuid().defaultRandom().primaryKey().notNull(),
	userId: uuid("user_id").notNull(),
	service: text().notNull(),
	accessToken: text("access_token").notNull(),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow().notNull(),
});
