import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '40s', target: 1000 }, // Ramp-up para 500 VUs
        { duration: '7s', target: 500 },  // Mantém 500 VUs
        { duration: '5s', target: 0 },   // Ramp-down
    ],
 };

export default function () {
    let url = 'http://localhost:8080/v1/user';  // Altere se necessário

    let payload = JSON.stringify({
        name: "Jefferson",
        year: 2024
    });

    let params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    let res = http.post(url, payload, params);

    // Verifica se o status da resposta é 200 ou 201 (Criado)
    check(res, {
        'status é 200 ou 201': (r) => r.status === 200 || r.status === 201,
        'resposta contém JSON': (r) => r.headers['Content-Type'] === 'application/json',
    });

    //console.log(`Status: ${res.status}, Resposta: ${res.body}`);

   // sleep(1); // Espera 1s antes de fazer outra requisição
}
