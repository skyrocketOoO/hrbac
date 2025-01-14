import http from 'k6/http';
import { check, sleep } from 'k6';

// Configuration for the test
export let options = {
    stages: [
        { duration: '10s', target: 10 }, // Ramp-up to 10 users over 10 seconds
        { duration: '30s', target: 10 }, // Stay at 10 users for 30 seconds
        { duration: '10s', target: 0 },  // Ramp-down to 0 users over 10 seconds
    ],
    vus: 1,         
};

export default function () {
  // First, Add 10,000 Permissions
  const count = 10000;
  for (let i = 0; i < count; i++) {
    addPermission();
  }
  for (let i = 0; i < count; i++) {
    checkPermission();
  }
}

// Test for the /checkPermission API
export function checkPermission() {
    // Set query parameters
    let url = 'http://localhost:8080/checkPermission';
    let params = {
      RoleName: 'Admin',
      PermissionName: 'Read',
      ObjectName: 'Document',
    };

    // Make the GET request
    let response = http.get(`${url}?RoleName=${params.RoleName}&PermissionName=${params.PermissionName}&ObjectName=${params.ObjectName}`);

    // Check for a successful response
    check(response, {
        'status is 200': (r) => r.status === 200,
        'response time is below 500ms': (r) => r.timings.duration < 500,
    });
}

// Test for the /addPermission API
export function addPermission() {
    let url = 'http://localhost:8080/addPermission';
    let payload = JSON.stringify({
      RoleName: 'Admin',
        PermissionName: 'Read',
        ObjectName: 'Document',
    });

    let params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // Make the POST request
    let response = http.post(url, payload, params);

    // Check for a successful response
    check(response, {
        'status is 200': (r) => r.status === 200,
        'response time is below 500ms': (r) => r.timings.duration < 500,
    });
}
