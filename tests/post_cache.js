import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 50,             
  duration: '30s',     
};

export default function () {
  let url = 'http://localhost:8080/cache';  

  
  let payload = JSON.stringify({
    id: "3",
    item: "Make bed",
    completed: false
  });

  
  let headers = {
    'Content-Type': 'application/json',
  };

  
  let res = http.post(url, payload, { headers: headers });

  
  check(res, {
    'status is 201': (r) => r.status === 201,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1); 
}