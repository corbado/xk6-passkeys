import http from 'k6/http';
import { check } from 'k6';
import passkeys from 'k6/x/passkeys';

export const options = {
    vus: 2,
    duration: "60s",
};

const baseUrl = 'http://localhost:8080';
const rp = passkeys.newRelyingParty('WebAuthn Demo', 'localhost', 'http://localhost:8080');

export default function () {
    const username = Math.random().toString(36).substring(2, 22);

    // Step 1: Start registration
    const startResponse = http.get(`${baseUrl}/register/start/${username}`, { tags: { name: 'register/start' } });
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
            tags: { name: 'register/finish' },
        }
    );

    check(finishResponse, {
        'registration finish status is 200': (r) => r.status === 200,
        'registration finish is successful': (r) => r.json('status') === 'Registration Success',
    });
}