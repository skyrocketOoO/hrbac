export const options = {
  scenarios: {
    average_load: {
      executor: 'constant-vus',
      vus: 50, // Expected average number of users
      duration: '5m', // Simulate normal conditions for 5 minutes
    },
  },
};
