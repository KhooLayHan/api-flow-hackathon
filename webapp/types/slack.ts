// export type SlackOAuthResponse {
//   ok: boolean;
//   authed_user: {
//     access_token: string;
//   };
//   error?: string;
//   access_token?: string;
//   scope?: string;
//   user_id?: string;
//   team_id?: string;
//   bot_user_id?: string;
//   bot_access_token?: string;
// }

export interface SlackOAuthResponse {
  ok: boolean;
  access_token: string;
  scope: string;
  token_type: string;
  authed_user?: {
    access_token: string;
  };
}

export interface SlackUser {
  id: string;
  email: string;
  name: string;
  team_id: string;
}
