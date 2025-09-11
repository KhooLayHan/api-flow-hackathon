export const customNodes = {
  githubTrigger: {
    type: 'githubTrigger',
    label: 'GitHub Push',
    data: {
      repository: 'https://github.com/KhooLayHan/api-flow-hackathon',
      branch: 'main',
    },
  },
  slackAction: {
    type: 'slackAction',
    label: 'Send Slack Message',
    data: {
      channel: '#general',
      message: 'A new commit was pushed!',
    },
  },
} as const;

export type CustomNodeKeys = keyof typeof customNodes;
