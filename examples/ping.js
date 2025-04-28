import http from 'k6/http';
import { success, failure } from './helper.js';

export const options = {
    vus: 2,
    duration: "30s",
};

const baseUrl = 'http://localhost:8080';

export default function () {
    const resp = http.get(`${baseUrl}/ping`);
    if (resp.status !== 200) {
        failure(`Request to ping failed with status ${resp.status}`);
    }

    success();
}