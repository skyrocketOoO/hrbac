export const options = {
  scenarios: {
    smoke: {
      executor: 'constant-vus',
      vus: 1, // Minimal load
      duration: '10s', // Short duration
    },
  },
};