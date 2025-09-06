<script setup>
const user = useSupabaseUser();
const supabase = useSupabaseClient();

async function handleLogout() {
  const { error } = await supabase.auth.signOut();
  if (error) {
    console.error(`Error logging out: ${error}`);
  }

  // Navigate to the login page after a successful sign-out
  await navigateTo('/login');
}
</script>

<template>
  <div>
    <!-- Global Header / Navigation Bar -->
    <header class="bg-white shadow-sm">
      <nav class="container mx-auto px-4 py-3 flex justify-between items-center">
        <!-- Logo Name -->
        <NuxtLink to="/" class="text-xl font-bold text-gray-800"> API Flow Hackathon </NuxtLink>

        <!-- Navigation Links & User Status -->
        <div class="flex items-center space-x-4">
          <!-- Show Login link if user is logged out -->
          <NuxtLink v-if="!user" to="/login" class="text-gray-600 hover:text-gray-800"> Login </NuxtLink>

          <!-- Show Dashboard link if user is logged in -->
          <NuxtLink v-if="user" to="/dashboard" class="text-gray-600 hover:text-gray-800"> Dashboard </NuxtLink>

          <!-- Show user email and Logout button if user is logged in -->
          <div v-if="user" class="flex items-center space-x-2">
            <span class="text-sm text-gray-500">{{ user.email }}</span>
            <button class="px-3 py-1 text-sm bg-red-500 text-white rounded hover:bg-red-600" @click="handleLogout">
              Logout
            </button>
          </div>
        </div>
      </nav>
    </header>

    <!-- Main Content Area -->
    <main class="container mx-auto px-4 py-8">
      <NuxtPage />
    </main>

    <!-- Global Footer -->
    <footer class="bg-gray-100 py-4">
      <div class="container mx-auto px-4 text-center">
        <p class="text-sm text-gray-500">© 2025 API Flow Hackathon</p>
      </div>
    </footer>
  </div>
</template>
