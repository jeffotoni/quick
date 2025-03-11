import http from 'k6/http';
import { check } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 500 },
        { duration: '5s', target: 500 }, 
        { duration: '2s', target: 0 },   
    ],
};

export default function () {

     let randomId = Math.floor(Math.random() * 1000); 
    let url = `http://localhost:8080/v1/user/${randomId}`;

     const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    let res = http.get(url, params, {
        timeout: '5s', // 5 second timeout
    });

    check(res, {
        'status is 200': (r) => r.status === 200,
        'no errors': (r) => !r.error,
    });
}