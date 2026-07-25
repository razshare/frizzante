export function keyof<T extends Record<string, unknown>>(key: keyof T): keyof T {
    return key
}