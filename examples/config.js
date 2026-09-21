// Defaults are a small local demonstration. RATE switches to ceremonies/second.
export function numberSetting(name, fallback, minimum = 1) {
  const value = Number(__ENV[name] === undefined ? fallback : __ENV[name]);
  if (!Number.isFinite(value) || value < minimum) {
    throw new Error(`${name} must be a number >= ${minimum}`);
  }
  return value;
}

function integerSetting(name, fallback) {
  const value = numberSetting(name, fallback);
  if (!Number.isInteger(value)) throw new Error(`${name} must be an integer`);
  return value;
}

export const baseUrl = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/$/, '');
export const rpId = __ENV.RP_ID || 'localhost';
export const origin = __ENV.ORIGIN || 'http://localhost:8080';
export const thinkTime = numberSetting('THINK_TIME', 0, 0);
const vus = integerSetting('VUS', 2);
const rate = __ENV.RATE === undefined ? null : integerSetting('RATE', 1);
const preAllocatedVUs = integerSetting('PRE_ALLOCATED_VUS', vus);
const maxVUs = integerSetting('MAX_VUS', preAllocatedVUs);
if (maxVUs < preAllocatedVUs) throw new Error('MAX_VUS must be >= PRE_ALLOCATED_VUS');
export const userCount = integerSetting('USER_COUNT', rate === null ? vus : maxVUs);
const duration = __ENV.DURATION || '30s';

export const options = {
  setupTimeout: __ENV.SETUP_TIMEOUT || '2m',
  // Omit the URL tag: randomly generated usernames otherwise create a series per user.
  systemTags: ['status', 'method', 'name', 'scenario', 'expected_response', 'error_code'],
  scenarios: {
    passkeys: rate === null ? {
      executor: 'constant-vus', vus, duration,
    } : {
      executor: 'constant-arrival-rate', rate, timeUnit: '1s', duration, preAllocatedVUs, maxVUs,
    },
  },
  thresholds: {
    checks: ['rate>0.99'],
    http_req_failed: ['rate<0.01'],
    ceremony_success: ['rate>0.99'],
    ceremony_duration: [`p(95)<${numberSetting('CEREMONY_P95_MS', 1000)}`],
    ...(rate === null ? {} : { dropped_iterations: ['count==0'] }),
  },
};
