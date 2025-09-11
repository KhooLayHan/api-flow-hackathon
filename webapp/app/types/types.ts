import type { Element, GraphNode } from '@vue-flow/core';

export interface GithubTriggerData {
  repository: string;
  branch: string;
}

export interface SlackActionData {
  channel: string;
  message: string;
}

export type CustomNodeData = GithubTriggerData | SlackActionData;

export type CustomElement = Element<CustomNodeData>;

export type CustomGraphNode = GraphNode<CustomNodeData>;

export interface WorkflowDefinition {
  elements: CustomElement;
}

export interface Workflow {
  id: number;
  name: string;
  githubRepo: string | null;
  definition: WorkflowDefinition;
  createdAt: string;
  updatedAt: string;
}

export interface NodeUpdatePayload {
  nodeId: string;
  data: CustomNodeData;
}
