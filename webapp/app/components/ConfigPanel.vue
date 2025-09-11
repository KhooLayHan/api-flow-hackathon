<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps({
  selectedNode: {
    type: Object,
    default: null,
  },
});

const emit = defineEmits(['update:node']);
const editableNodeData = ref(null);

watch(
  () => props.selectedNode,
  newNode => {
    if (newNode) {
      // Deep copy to avoid mutating the original object
      editableNodeData.value = JSON.parse(JSON.stringify(newNode));
    } else {
      editableNodeData.value = null;
    }
  },
  { immediate: true, deep: true },
);

function onDataChange() {
  if (props.selectedNode) {
    emit('update:node', {
      nodeId: props.selectedNode.id,
      data: editableNodeData.value,
    });
  }
}
</script>

<template>
  <div v-if="selectedNode && editableNodeData" class="w-80 bg-gray-50 p-4 border-l">
    <h3 class="font-bold text-lg mb-4">Configure Node</h3>
    <p class="text-sm text-gray-500 mb-4">ID: {{ selectedNode.id }}</p>

    <!-- GitHub Trigger Configuration -->
    <div v-if="selectedNode.type === 'githubTrigger'">
      <label class="block text-sm font-medium">Repository</label>
      <input
        v-model="editableNodeData.repository"
        type="text"
        class="p-2 border rounded w-full"
        @input="onDataChange"
      />
    </div>

    <!-- Slack Action Configuration -->
    <div v-if="selectedNode.type === 'slackAction'">
      <label class="block text-sm font-medium">Channel</label>
      <input
        v-model="editableNodeData.channel"
        type="text"
        class="p-2 border rounded w-full mb-2"
        @input="onDataChange"
      />
      <label class="block text-sm font-medium">Message Template</label>
      <textarea v-model="editableNodeData.message" class="p-2 border rounded w-full" rows="4" @input="onDataChange" />
      <p class="text-xs text-gray-400 mt-1">Use variables {commit.message}, {pusher.name}, and {repo.name}.</p>
    </div>
  </div>
</template>
