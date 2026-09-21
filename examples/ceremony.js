import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import passkeys from 'k6/x/passkeys';
import {
  baseUrl, rpId, origin, thinkTime,
} from './config.js';
import { randomString } from './helper.js';

const rp = passkeys.newRelyingParty('WebAuthn Demo', rpId, origin);
const ceremonySuccess = new Rate('ceremony_success');
const ceremonyDuration = new Trend('ceremony_duration', true);

function requireOK(response, endpoint) {
  if (response.status !== 200) {
    throw new Error(`${endpoint} returned HTTP ${response.status}`);
  }
}

export function register() {
  const username = randomString(20);
  const start = http.get(`${baseUrl}/register/start/${username}`, {
    tags: { name: 'register/start' },
  });
  requireOK(start, 'register/start');
  if (thinkTime) sleep(thinkTime);
  const credential = passkeys.newCredential();
  const response = passkeys.createAttestationResponse(rp, credential, start.body);
  const finish = http.post(`${baseUrl}/register/finish/${username}`, response, {
    headers: { 'Content-Type': 'application/json' },
    tags: { name: 'register/finish' },
  });
  requireOK(finish, 'register/finish');
  // Explicit JSON preserves the credential's custom key serialization across setup/VUs.
  // This demo uses the username's UTF-8 bytes as the WebAuthn user ID.
  return { username, userHandle: username, credential: JSON.stringify(credential) };
}

export function login(user) {
  const start = http.get(`${baseUrl}/login/start/${user.username}`, {
    tags: { name: 'login/start' },
  });
  requireOK(start, 'login/start');
  if (thinkTime) sleep(thinkTime);
  const credential = JSON.parse(user.credential);
  const response = passkeys.createAssertionResponse(rp, credential, user.userHandle, start.body);
  const finish = http.post(`${baseUrl}/login/finish/${user.username}`, response, {
    headers: { 'Content-Type': 'application/json' },
    tags: { name: 'login/finish' },
  });
  requireOK(finish, 'login/finish');
}

export function measure(run) {
  const started = Date.now();
  let completed = false;
  try {
    run();
    completed = true;
  } finally {
    // Count exceptions as failures too, including JSON and virtual-authenticator errors.
    ceremonySuccess.add(completed);
    ceremonyDuration.add(Date.now() - started);
    check(completed, { success: (value) => value });
  }
}
