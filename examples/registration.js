import http from 'k6/http';
import passkeys from 'k6/x/passkeys';
import { success, failure, randomString } from './helper.js';

export const options = {
  vus: 2,
  duration: '30s',
};

const baseUrl = 'http://localhost:8080';
const rp = passkeys.newRelyingParty('WebAuthn Demo', 'localhost', 'http://localhost:8080');

export default function () {
  const username = randomString(20);

  // Step 1: Start registration
  const startResponse = http.get(`${baseUrl}/register/start/${username}`, { tags: { name: 'start' } });
  if (startResponse.status !== 200) {
    failure(`Request to register/start failed with status ${startResponse.status} (body: ${startResponse.body})`);
  }

  // Step 2: Create attestation response (simulate the client
  // side and call to navigator.credentials.create())
  const credential = passkeys.newCredential();
  const attestationResponse = passkeys.createAttestationResponse(
    rp,
    credential,
    JSON.stringify(startResponse.json()),
  );

  // Step 3: Finish registration
  const finishResponse = http.post(
    `${baseUrl}/register/finish/${username}`,
    attestationResponse,
    {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: 'finish' },
    },
  );
  if (finishResponse.status !== 200) {
    failure(`Request to register/finish failed with status ${finishResponse.status} (body: ${finishResponse.body})`);
  }

  success();
}
