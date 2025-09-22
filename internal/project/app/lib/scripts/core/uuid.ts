/**
 * Create a [universally unique identifier](https://en.wikipedia.org/wiki/Universally_unique_identifier).
 * @param short if set to true, the resulting string will be 8 characters long instead of 36.
 * @returns
 */
export function uuid(short: boolean = false): string {
    let time = new Date().getTime()
    const blueprint = short ? "xyxxyxyx" : "xxxxxxxx-xxxx-yxxx-yxxx-xxxxxxxxxxxx"
    const result = blueprint.replace(/[xy]/g, function check(char) {
        const value = (time + Math.random() * 16) % 16 | 0
        time = Math.floor(time / 16)
        return (char == "x" ? value : (value & 0x3) | 0x8).toString(16)
    })
    return result
}
