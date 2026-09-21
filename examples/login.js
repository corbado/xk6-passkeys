import exec from 'k6/execution';
import { register, login, measure } from './ceremony.js';
import { userCount } from './config.js';

export { options } from './config.js';

export function setup() {
  return Array.from({ length: userCount }, () => register());
}

export default function (users) {
  // One identity per VU when USER_COUNT covers all VUs; smaller pools intentionally share.
  const user = users[(exec.vu.idInTest - 1) % users.length];
  measure(() => login(user));
}
