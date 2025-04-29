import { check, fail } from "k6";

/**
 * Generates a random string of a given length
 * @param {number} length - The length of the string to generate
 * @returns {string} - A random string of the specified length
 */
export function randomString(length) {
    return Array.from({ length: length }, () => {
        const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
        return chars.charAt(Math.floor(Math.random() * chars.length));
    }).join('');
}

/**
 * Success helper that calls k6 check with true to count for an successful iteration
 */
export function success() {
    check(true, { success: (r) => r === true });
}

/**
 * Failure helper that calls k6 check with false to count for an failed iteration
 * @param {string} message - The message to log
 */
export function failure(message) {
    check(false, { success: (r) => r === true });
    fail(message);
}