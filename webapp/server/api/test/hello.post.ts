import { defineEventHandler } from "h3";
import { workflowQueue } from "../../utils/queue";

export default defineEventHandler(async (event) => {
  console.log("Received request to /api/test/hello");

  const jobName = "test-message";
  const jobData = {
    message: `Hello from Nuxt! The time is ${new Date().toLocaleTimeString()}`,
    sentAt: new Date().toISOString(),
  };

  // Adding the queue job;
  await workflowQueue.add(jobName, jobData);
  console.log(`Job added to queue: ${jobName}, ${jobData.message}`);

  return {
    status: "ok",
    message: "Job has been successfully added to the queue.",
  };
});
