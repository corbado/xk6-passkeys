import { register, measure } from './ceremony.js';

export { options } from './config.js';

export default function () {
  measure(register);
}
