
import { check, fail } from "k6";

export function success() {
    check(true, { success: (r) => r === true });
}

export function failure(message) {
    check(false, { success: (r) => r === true });
    fail(message);
}