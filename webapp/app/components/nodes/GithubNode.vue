<script setup>
import { Handle, Position } from '@vue-flow/core';
</script>

<template>
  <div class="p-4 border border-gray-300 rounded-md bg-white">
    <!-- Icon and Title -->
    <div class="flex items-center space-x-2">
      <Icon name="mdi:github" class="h-6 w-6" />
      <div class="font-bold">Github Push Trigger</div>
    </div>
    <!-- Output handle -->
    <Handle type="source" :position="Position.Right" />
  </div>
</template>
