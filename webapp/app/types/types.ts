import type { Element } from '@vue-flow/core';

export interface GithubTriggerData {
  repository: string;
  branch: string;
}

export interface SlackActionData {
  channel: string;
  message: string;
}

export type CustomNodeData = GithubTriggerData | SlackActionData;

export interface WorkflowDefinition {
  elements: Element<CustomNodeData>[];
}

export interface NodeUpdatePayload {
  nodeId: string;
  data: CustomNodeData;
}
