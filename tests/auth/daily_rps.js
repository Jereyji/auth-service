import http from 'k6/http';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '3m', target: 300 },
    { duration: '3m', target: 100 },
    { duration: '2m', target: 0 },  
  ],
};

export default function () {
  let url = 'http://localhost:8080/auth/login';
  let payload = JSON.stringify({
    email: 'test@test.com',
    password: 'test1234',
  });

  let params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  let res = http.post(url, payload, params);

  check(res, { 'status is 200': (r) => r.status === 200 });
}
