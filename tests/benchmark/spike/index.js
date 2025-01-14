export const options = {
  scenarios: {
    spike: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '10s', target: 200 }, // Sudden spike to 200 users
        { duration: '1m', target: 200 },  // Maintain spike load
        { duration: '10s', target: 0 },   // Drop back to 0 users
      ],
    },
  },
};