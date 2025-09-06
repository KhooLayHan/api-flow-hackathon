<script setup>
const user = useSupabaseUser();
const router = useRouter();

watch(
  user,
  () => {
    if (user.value) {
      return router.push('/');
    }
  },
  { immediate: 'true' },
);
</script>

<template>
  <div>
    <h1>Please wait while we log you in...</h1>
  </div>
</template>
