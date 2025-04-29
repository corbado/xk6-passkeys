import http from 'k6/http';
import passkeys from 'k6/x/passkeys';
import { success, failure, randomString } from './helper.js';

export const options = {
  vus: 2,
  duration: '30s',
};

const baseUrl = 'http://localhost:8080';
const rp = passkeys.newRelyingParty('WebAuthn Demo', 'localhost', 'http://localhost:8080');

// Setup function to create a single test user
export function setup() {
  const username = randomString(20);

  // Step 1: Start registration
  const startResponse = http.get(`${baseUrl}/register/start/${username}`);
  if (startResponse.status !== 200) {
    throw new Error(`Request to register/start failed with status ${startResponse.status} (body: ${startResponse.body})`);
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
    },
  );

  if (finishResponse.status !== 200) {
    throw new Error(`Request to register/finish failed with status ${finishResponse.status} (body: ${finishResponse.body})`);
  }

  // We need to stringify the credential to avoid the
  // "invalid credential" error (did not find the root cause yet)
  return { username, credential: JSON.stringify(credential) };
}

export default function (data) {
  const { username } = data;
  const credential = JSON.parse(data.credential);

  // Step 1: Start login
  const startResponse = http.get(`${baseUrl}/login/start/${username}`, { tags: { name: 'start' } });
  if (startResponse.status !== 200) {
    failure(`Request to login/start failed with status ${startResponse.status} (body: ${startResponse.body})`);
  }

  // Step 2: Create assertion response (simulate the client
  // side and call to navigator.credentials.get())
  const assertionResponse = passkeys.createAssertionResponse(
    rp,
    credential,
    username,
    JSON.stringify(startResponse.json()),
  );

  // Step 3: Finish login
  const finishResponse = http.post(
    `${baseUrl}/login/finish/${username}`,
    assertionResponse,
    {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: 'finish' },
    },
  );
  if (finishResponse.status !== 200) {
    failure(`Request to login/finish failed with status ${finishResponse.status} (body: ${finishResponse.body})`);
  }

  success();
}
