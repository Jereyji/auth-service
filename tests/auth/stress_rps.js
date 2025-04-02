import http from 'k6/http';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 1000 }, // Разгон до 1000 RPS за 10 секунд
    { duration: '30s', target: 1000 }, // Держим нагрузку
    { duration: '10s', target: 0 }, // Снижаем нагрузку
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
