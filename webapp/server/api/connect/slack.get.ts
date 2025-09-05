import { defineEventHandler, sendRedirect } from 'h3';

export default defineEventHandler(async event => {
  const clientId = process.env.SLACK_CLIENT_ID;
  if (!clientId) {
    throw new Error(`SLACK_CLIENT_ID is not set: ${clientId}`);
  }

  // Permission (scopes)
  const scopes = 'chat:write,channels:read';

  // Slack Authorization URL
  const authUrl = new URL('https://slack.com/oauth/v2/authorize');

  authUrl.searchParams.set('client_id', clientId);
  authUrl.searchParams.set('scope', scopes);
  authUrl.searchParams.set('redirect_uri', 'http://localhost:3000/api/auth/callback/slack');

  await sendRedirect(event, authUrl.toString());
});
