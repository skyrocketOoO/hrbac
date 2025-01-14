import http from 'k6/http';
import { check, group } from 'k6';
import { Trend } from 'k6/metrics';

// Custom metrics for tracking API performance
const pingDuration = new Trend('ping_duration');
const healthyDuration = new Trend('healthy_duration');

const BASE_URL = 'http://localhost:8080/v1'; // Replace with your server's base URL

export let options = {
  vus: 1,       
  iterations: 1, 
};

export default function () {
  group('Ping Endpoint', () => {
    const start = Date.now();
    let res = http.get(`${BASE_URL}/ping`);
    pingDuration.add(Date.now() - start);


    check(res, {
      'Ping: status is 200': (r) => r.status === 200,
    });
  });

  group('Healthy Endpoint', () => {
    const start = Date.now();
    let res = http.get(`${BASE_URL}/healthy`);
    healthyDuration.add(Date.now() - start);

    check(res, {
      'Healthy: status is 200': (r) => r.status === 200,
      'Healthy: response time is < 200ms': (r) => r.timings.duration < 200,
    });
  });
}
