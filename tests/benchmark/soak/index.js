export const options = {
  scenarios: {
    soak: {
      executor: 'constant-vus',
      vus: 20, // Moderate load
      duration: '12h', // Extended duration
    },
  },
};