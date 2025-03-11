import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '10s', target: 250 }, // Scales to 50 users in 30 seconds
        //{ duration: '20s', target: 1000 }, // Keeps 100 users for 1 minute
        { duration: '15s', target: 0 },  // Reduces to 0 users in 30 seconds
    ],
    http2: true,
    noConnectionReuse: false, // Reuse HTTP/2 connections
    insecureSkipTLSVerify: true, // Skip TLS certificate verification
    batchPerHost: 100,
};

export default function () {
    let randomId = Math.floor(Math.random() * 1000); // Generates an ID between 0 and 999
    let url = `https://localhost:443/v1/user/${randomId}`;

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    let res = http.get(url, params, {
        timeout: '5s', // 5 second timeout
    });

    // Check if the request was successful
    check(res, {
        'status is 200': (r) => r.status === 200,
        'protocol is HTTP/2': (r) => r.proto === 'HTTP/2.0',
        'no errors': (r) => !r.error,
    });

    // Logs for debugging
    // console.log(`debug..: ${res.body}`)
    // console.log(`Requesting user with ID: ${randomId}`);
    // console.log(`Response status: ${res.status}, Protocol: ${res.proto}`);
    // sleep(1); // 1 second interval between requests
}