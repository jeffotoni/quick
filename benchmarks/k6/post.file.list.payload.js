import http from 'k6/http';
import { check, sleep } from 'k6';

// Load the JSON from the environment variable
const payloadData = open('./data_1k_list.json');

// K6 configuration
export let options = {
    stages: [
        { duration: '40s', target: 1000 }, // Ramp-up para 500 VUs
        { duration: '7s', target: 500 },  // Mantém 500 VUs
        { duration: '5s', target: 0 },   // Ramp-down
    ],
 };


export default function () {
let url = 'http://localhost:8080/v1/user';

// Always use the same list for sending
// let payload = JSON.stringify(payloadData);

let params = {
headers: { 'Content-Type': 'application/json' },
};

let res = http.post(url, payloadData, params);

// Check if the response is correct
check(res, {
'status is 200 or 201': (r) => r.status === 200 || r.status === 201,
'response contains JSON': (r) => r.headers['Content-Type'] === 'application/json',
});

// console.log(`Status: ${res.status}, Response: ${res.body}`);
// sleep(1); // Wait 1 second before repeating the request

}