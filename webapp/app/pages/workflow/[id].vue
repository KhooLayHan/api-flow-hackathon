<script setup lang="ts">
import { ref, watchEffect } from 'vue';
import { VueFlow } from '@vue-flow/core';
import { useRoute } from 'vue-router';

import '@vue-flow/core/dist/style.css';
import '@vue-flow/core/dist/theme-default.css';

// 1. Setup
const route = useRoute();
const workflowId = Number(route.params.id);
const elements = ref([]);

// const elements = ref([
//   // Hardcoded
//   { id: '1', type: 'input', label: 'Trigger', position: { x: 100, y: 100 } },
//   { id: '2', label: 'Log Message', position: { x: 300, y: 100 } },
//   { id: 'e1-2', source: '1', target: '2' },
// ]);

// 2. Loading the data with useAsyncData
const {
  data: workflowData,
  pending,
  error,
} = await useAsyncData(`workflow-${workflowId}`, () => $fetch(`/api/workflows/${workflowId}`));

watchEffect(() => {
  if (!workflowData.value) return;

  console.log(workflowData.value);

  // Process workflowData.value and update elements accordingly
  if (workflowData.value?.workflow?.definition) {
    elements.value = workflowData.value.workflow.definition;

    console.log(elements.value);
  }
});

// 3. Save the data
async function saveWorkflow() {
  console.log('Saving workflow...');

  try {
    const response = await $fetch(`/api/workflows/${workflowId}`, {
      method: 'GET',
      body: {
        name: workflowData.value?.workflow?.name || 'My Workflow',
        definition: elements.value,
      },
    });

    console.log(`Workflow saved successfully, ${response}!`);
  } catch (error) {
    console.error(`Error saving workflow: ${error}`);
  }
}
</script>

<template>
  <!-- Top bar with save button -->
  <div class="flex justify-between items-center p-4">
    <h1 class="text-2xl font-bold">Workflow</h1>
    <button
      class="px-8 py-16 bg-color-#007bff color-white border-0 border-radius-4px cursor-pointer"
      @click="saveWorkflow"
    >
      Save
    </button>
  </div>

  <!-- Loading indicator using the 'pending' state from useAsyncData -->
  <div v-if="pending" class="flex justify-center items-center h-full">
    <span>Loading workflow...</span>
  </div>

  <!-- Error message -->
  <div v-else-if="error" class="flex justify-center items-center h-full">
    <button>Error loading workflow: {{ error }}</button>
  </div>

  <!-- Main Vue canvas -->
  <div v-else class="h-full w-full">
    <VueFlow v-model="elements">
      <!-- TODO: Implement workflow controls, minimaps, etc. -->
    </VueFlow>
  </div>
</template>
