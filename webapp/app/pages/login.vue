<script setup>
const supabase = useSupabaseClient();

async function loginWithGitHub() {
  const { error } = await supabase.auth.signInWithOAuth({
    provider: 'github',
    options: {
      // redirectTo: window.location.origin + '/login/callback',
      redirectTo: 'http://localhost:3000/',
    },
  });

  if (error) {
    console.error(`Error logging in with GitHub: ${error}.`);
  }
}
</script>

<template>
  <div>
    <h1>Welcome to API Flow Hackathon</h1>
    <p>Please login to continue.</p>
    <button @click="loginWithGitHub">Login with GitHub</button>
  </div>
</template>
