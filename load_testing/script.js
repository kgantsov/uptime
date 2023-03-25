import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 200 },
        { duration: '40s', target: 100 },
        { duration: '20s', target: 10 },
    ],
};

export default function () {
    const url = `${__ENV.UPTIME_HOST}/API/v1/services`;

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${__ENV.UPTIME_TOKEN}`,
        },
    };

    const res = http.get(url, params);

    check(res, { 'status was 200': (r) => r.status == 200 });

    sleep(1)
}
