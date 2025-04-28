
import { check, fail } from "k6";

export function randomString(length) {
    return Array.from({ length: length }, () => {
        const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
        return chars.charAt(Math.floor(Math.random() * chars.length));
    }).join('');
}

export function success() {
    check(true, { success: (r) => r === true });
}

export function failure(message) {
    check(false, { success: (r) => r === true });
    fail(message);
}