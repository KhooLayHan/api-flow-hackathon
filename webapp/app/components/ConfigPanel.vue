<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import type { CustomGraphNode, NodeUpdatePayload, GithubTriggerData, SlackActionData } from '@/types/types';

const props = defineProps<{
  selectedNode: CustomGraphNode | null;
}>();

const emit = defineEmits<{
  (e: 'update:node', payload: NodeUpdatePayload): void;
}>();

const editableNodeData = ref<GithubTriggerData | SlackActionData | null>(null);

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
  if (props.selectedNode && editableNodeData.value) {
    emit('update:node', {
      nodeId: props.selectedNode.id,
      data: editableNodeData.value,
    });
  }
}

const isGithubTrigger = computed(() => props.selectedNode?.type === 'githubTrigger' && editableNodeData.value);
const isSlackAction = computed(() => props.selectedNode?.type === 'slackAction' && editableNodeData.value);
</script>

<template>
  <div v-if="selectedNode && editableNodeData" class="w-80 bg-gray-50 p-4 border-l">
    <h3 class="font-bold text-lg mb-4">Configure Node</h3>
    <p class="text-sm text-gray-500 mb-4">ID: {{ selectedNode.id }}</p>

    <!-- GitHub Trigger Configuration -->
    <div v-if="isGithubTrigger">
      <label class="block text-sm font-medium">Repository</label>
      <input
        v-model="(editableNodeData as GithubTriggerData).repository"
        type="text"
        class="p-2 border rounded w-full"
        @input="onDataChange"
      />
    </div>

    <!-- Slack Action Configuration -->
    <div v-if="isSlackAction">
      <label class="block text-sm font-medium">Channel</label>
      <input
        v-model="(editableNodeData as SlackActionData).channel"
        type="text"
        class="p-2 border rounded w-full mb-2"
        @input="onDataChange"
      />
      <label class="block text-sm font-medium">Message Template</label>
      <textarea
        v-model="(editableNodeData as SlackActionData).message"
        class="p-2 border rounded w-full"
        rows="4"
        @input="onDataChange"
      />
      <p class="text-xs text-gray-400 mt-1">Use variables {commit.message}, {pusher.name}, and {repo.name}.</p>
    </div>
  </div>
</template>
