import http from 'k6/http';
import { check } from 'k6';
import passkeys from 'k6/x/passkeys';

export const options = {
    stages: [
        { duration: '30s', target: 10 }, // Ramp up to 10 users
        { duration: '1m', target: 10 },  // Stay at 10 users
        { duration: '30s', target: 0 },  // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
        http_req_failed: ['rate<0.01'],   // Less than 1% of requests should fail
    },
};

const baseUrl = 'http://localhost:8080';
const rp = passkeys.newRelyingParty('WebAuthn Demo', 'localhost', 'http://localhost:8080');

export default function () {
    // Generate a unique username for each virtual user
    const username = `user_${__VU}_${Date.now()}`;

    // Step 1: Start registration
    const startResponse = http.get(`${baseUrl}/register/start/${username}`);
    check(startResponse, {
        'registration start status is 200': (r) => r.status === 200,
        'registration start has options': (r) => r.json() !== null,
    });

    if (startResponse.status !== 200) {
        return;
    }

    // Step 2: Create attestation response
    const credential = passkeys.newCredential();
    const attestationResponse = passkeys.createAttestationResponse(
        rp,
        credential,
        JSON.stringify(startResponse.json())
    );

    // Step 3: Finish registration
    const finishResponse = http.post(
        `${baseUrl}/register/finish/${username}`,
        attestationResponse,
        {
            headers: { 'Content-Type': 'application/json' },
        }
    );

    check(finishResponse, {
        'registration finish status is 200': (r) => r.status === 200,
        'registration finish is successful': (r) => r.json('status') === 'Registration Success',
    });
}