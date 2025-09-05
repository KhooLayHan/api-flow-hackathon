// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxt/eslint', '@nuxt/icon', '@nuxt/image', '@nuxtjs/supabase'],

  supabase: {
    clientOptions: {
      auth: {
        flowType: 'pkce', // Recommended for OAuth
        autoRefreshToken: true,
        detectSessionInUrl: true,
        // redirectURI: 'http://localhost:3000/auth/callback',
        // scopes: 'email,profile'
      },
    },
    // Handles the user being redirected to the app after logging in
    redirectOptions: {
      login: '/login', // Redirect here if user is not logged in
      callback: '/confirm', // Redirect here after successful login
      // error: '/error'
    },
  },

  runtimeConfig: {
    public: {
      supabase: {
        url: process.env.SUPABASE_URL,
        key: process.env.SUPABASE_KEY,
      },
    },
  },
});
