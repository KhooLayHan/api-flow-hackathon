<script setup lang="ts">
import { ref, watchEffect } from 'vue';
import { VueFlow, useVueFlow, type GraphNode, type Element } from '@vue-flow/core';
import { useRoute } from 'vue-router';

import { customNodes, type CustomNodeKeys } from '@/utils/nodes';
import GithubNode from '~/components/nodes/GithubNode.vue';
import SlackNode from '~/components/nodes/SlackNode.vue';
import ConfigPanel from '@/components/ConfigPanel.vue';

import '@vue-flow/core/dist/style.css';
import '@vue-flow/core/dist/theme-default.css';

import type { NodeUpdatePayload } from '@/types';

// 1. Setup
const route = useRoute();
const workflowId = Number(route.params.id);

// Core Vue Flow instance for programmatic access
const { onPaneClick, onNodeClick, addNodes, screenToFlowCoordinate } = useVueFlow();

const elements = ref<Element[]>([]);
const selectedNode = ref<GraphNode | null>(null);
const workflowName = ref('Loading...');

const isSaving = ref(false);

// const elements = ref([
//   // Hardcoded
//   { id: '1', type: 'input', label: 'Trigger', position: { x: 100, y: 100 } },
//   { id: '2', label: 'Log Message', position: { x: 300, y: 100 } },
//   { id: 'e1-2', source: '1', target: '2' },
// ]);

// 2. Fetch and loads the data.
const {
  data: workflowData,
  pending,
  error,
} = await useAsyncData(`workflow-${workflowId}`, () => $fetch(`/api/workflows/${workflowId}`));

watchEffect(() => {
  if (!workflowData.value) return;

  console.log(workflowData.value);

  if (workflowData.value?.workflow) {
    const workflow = workflowData.value.workflow;
    elements.value = workflow.definition.elements || []; // Use saved elements of empty array
    workflowName.value = workflow.name;

    console.log(elements.value);
  }
});

// Process workflowData.value and update elements accordingly
// if (workflowData.value?.workflow?.definition) {
//   elements.value = workflowData.value.workflow.definition || []; // Use saved elements of empty array

//   console.log(elements.value);
// }
//
// 3. Event handlers & actions

// Node selection for ConfigPanel
onNodeClick(event => {
  // Set the selectedNode to show the config panel
  selectedNode.value = event.node;
});

onPaneClick(() => {
  // Deselects the node when clicking on the background pane.
  selectedNode.value = null;
});

// Adding new nodes
function onAddNode(nodeType: CustomNodeKeys) {
  const newNodeTemplate = customNodes[nodeType];

  const newNode = {
    id: `node-${Date.now()}`,
    type: newNodeTemplate.type,
    label: newNodeTemplate.label,
    position: screenToFlowCoordinate({ x: 150, y: 150 }),
    data: { ...newNodeTemplate.data },
  };

  addNodes([newNode]);
}

// 3. Save the data
async function handleSaveWorkflow() {
  if (!workflowId) return;

  isSaving.value = true;
  console.log('Saving workflow...');

  // Find the GitHub trigger node, if it exists
  const githubTriggerNode = elements.value.find(el => el.type === 'githubTrigger');
  const repoName = githubTriggerNode?.data?.repository;

  const { data, error } = await useFetch(`/api/workflows/${workflowId}/update`, {
    method: 'PUT',
    body: {
      name: workflowName.value,
      definition: { elements: elements.value }, // Save the elements under a key
      githubRepo: repoName, // Pass the githubRepo if a githubTrigger node exists
    },
    immediate: false,
  });

  isSaving.value = false;

  if (error.value) {
    console.error(`Error saving workflow: ${error.value}`);
    return;
  }

  console.log(`Workflow saved successfully, ${data.value}!`);
  alert('Workflow saved!');
}

async function handleTriggerWorkflow() {
  try {
    const response = await $fetch(`/api/workflows/${workflowId}/trigger`, {
      method: 'POST',
      // body: {
      //   name: workflowData.value || 'My Workflow',
      //   definition: {
      //     elements: elements.value,
      //   },
      //   githubRepo: elements.value.find(el => el.type === 'githubTrigger')?.data.repository,
      // },
    });

    console.log(`Workflow triggered successfully! Job ID: ${response.jobId}`);
  } catch (error) {
    console.error(`Error triggering workflow: ${error}`);
  }
}

function handleNodeUpdate(updatePayload: NodeUpdatePayload) {
  const nodeToUpdate = elements.value.find(el => el.id === updatePayload.nodeId);

  if (nodeToUpdate) {
    nodeToUpdate.data = updatePayload.data;
    // elements.value = elements.value.map(el => el.id === nodeToUpdate.id ? nodeToUpdate : el);
  }
}
</script>

<template>
  <div class="flex h-screen w-screen bg-gray-200">
    <!-- Left sidebar: Node Palette -->
    <div class="w-64 bg-white p-4 border-r flex flex-col">
      <h2 class="text-lg font-bold mb-4">API Flow Editor</h2>
      <input v-model="workflowName" type="text" class="p-2 border rounded w-full mb-4" />
      <hr class="mb-4" />
      <h3 class="font-semibold mb-2 text-gray-700">Add Nodes</h3>
      >
      <button
        class="p-2 border rounded bg-gray-50 hover:bg-gray-100 w-full mb-2 text-left"
        @click="onAddNode('githubTrigger')"
      >
        GitHub Trigger
      </button>
      <button
        class="p-2 border rounded bg-gray-50 hover:bg-gray-100 w-full text-left"
        @click="onAddNode('slackAction')"
      >
        Slack Action
      </button>

      <div class="mt-auto">
        <button
          class="p-2 rounded bg-blue-500 hover:bg-blue-600 text-white w-full mb-2 disabled:bg-blue-300"
          :disabled="isSaving"
          @click="handleSaveWorkflow"
        >
          {{ isSaving ? 'Saving...' : 'Save' }}
        </button>
        <!-- <button class="p-2 rounded bg-blue-500 hover:bg-blue-600 text-white w-full mb-2 @click="handleSaveWorkflow">Run Manually</button> -->
      </div>
    </div>
    <!-- Top bar with save button -->
    <!-- <div class="flex justify-between items-center p-4">
    <h1 class="text-2xl font-bold">Workflow</h1>
    <button
      class="px-8 py-16 bg-color-#007bff color-white border-0 border-radius-4px cursor-pointer"
      @click="saveWorkflow"
    >
      Save
    </button>
  </div> -->

    <!-- Loading indicator using the 'pending' state from useAsyncData -->
    <!-- <div v-if="pending" class="flex justify-center items-center h-full">
    <span>Loading workflow...</span>
  </div> -->

    <!-- Error message -->
    <!-- <div v-else-if="error" class="flex justify-center items-center h-full">
    <button>Error loading workflow: {{ error }}</button>
  </div> -->

    <!-- Main Vue canvas -->
    <!-- <div v-else class="h-full w-full">
    <VueFlow v-model="elements">
      // TODO: Implement workflow controls, minimaps, etc. -->
    <!-- </VueFlow> -->

    <!-- Main canvas -->
    <div class="flex-1 relative">
      <div v-if="pending" class="absolute inset-0 flex items-center justify-center bg-white/50 z-20">Loading...</div>
      <div v-else-if="error" class="absolute inset-0 flex items-center justify-center bg-red-100 text-red-700">
        Error loading workflows: {{ error.message }}
      </div>
      <VueFlow v-model="elements" class="bg-white">
        <!-- Register the custom components for rendering. -->
        <template #node-githubTrigger="props">
          <GithubNode v-bind="props" />
        </template>
        <template #node-slackAction="props">
          <SlackNode v-bind="props" />
        </template>
      </VueFlow>
    </div>
    <!-- Right sidebar: Configuration Panel -->
    <ConfigPanel :selected-node="selectedNode" @update:node="handleNodeUpdate" />
  </div>
</template>
