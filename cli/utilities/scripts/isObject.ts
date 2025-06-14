export function isObject(target:unknown) {
    return (!!target) && (target.constructor === Object)
}