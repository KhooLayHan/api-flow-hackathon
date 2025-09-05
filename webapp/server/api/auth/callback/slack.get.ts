import { defineEventHandler, getQuery, sendRedirect } from 'h3';
import { z } from 'zod';
import { db } from '~~/server/utils/db';
import { userCredentials } from '~~/server/db/schema';
import { serverSupabaseUser } from '#supabase/server';

import type { SlackOAuthResponse } from '~~/types/slack';

const SlackCallbackQuerySchema = z.object({
  code: z.string().min(1).max(1000),
  state: z.string().min(1).max(1000),
});

export default defineEventHandler(async event => {
  // 1. Get the current logged-in user from the session
  const user = await serverSupabaseUser(event);
  if (!user) {
    throw createError({ statusCode: 401, message: 'You must be logged in to connect a service.' });
  }

  // 2. Validate the query parameters from Slack
  const query = getQuery(event);

  const validation = SlackCallbackQuerySchema.safeParse(query);
  if (!validation.success) {
    throw createError({ statusCode: 400, message: 'Invalid query parameters from Slack.' });
  }

  const { code } = validation.data;

  // 3. Exchange the code for an access token
  try {
    const clientId = process.env.SLACK_CLIENT_ID;
    if (!clientId) {
      throw createError({ statusCode: 500, message: 'SLACK_CLIENT_ID is not set.' });
    }

    const clientSecret = process.env.SLACK_CLIENT_SECRET;
    if (!clientSecret) {
      throw createError({ statusCode: 500, message: 'SLACK_CLIENT_SECRET is not set.' });
    }

    const response: SlackOAuthResponse = await $fetch('https://slack.com/api/oauth.v2.access', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: new URLSearchParams({
        client_id: clientId,
        client_secret: clientSecret,
        code: code,
        redirect_uri: 'http://localhost:3000/api/auth/callback/slack',
      }).toString(),
    });

    if (!response.ok || !response.authed_user?.access_token) {
      console.error('Slack OAuth error:', response);
      throw new Error('Slack did not respond with a valid access token.');
    }

    const accessToken = response.authed_user.access_token;

    // 4. Save the access token securely to the database
    await db
      .insert(userCredentials)
      .values({
        userId: user.id,
        service: 'slack',
        accessToken: accessToken,
      })
      .onConflictDoUpdate({
        target: [userCredentials.userId, userCredentials.service],
        set: {
          accessToken: accessToken,
        },
      });

    // 5. Redirect the user back to a success page
    return sendRedirect(event, '/dashboard?connected=slack');
  } catch (error) {
    console.error(`Slack OAuth token exchange error: ${error}`);
    return sendRedirect(event, '/dashboard?error=slack_connection_failed');
  }
});
