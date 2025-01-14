export const options = {
  scenarios: {
    stress: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '1m', target: 100 }, // Ramp up to 100 users
        { duration: '3m', target: 200 }, // Peak load with 200 users
        { duration: '1m', target: 0 },   // Ramp down to 0 users
      ],
    },
  },
};