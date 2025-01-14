export const options = {
  scenarios: {
    breakpoint: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '1m', target: 50 },  // Ramp up to 50 users
        { duration: '1m', target: 100 }, // Ramp up to 100 users
        { duration: '1m', target: 200 }, // Ramp up to 200 users
        { duration: '1m', target: 300 }, // Gradually increase further
      ],
    },
  },
};