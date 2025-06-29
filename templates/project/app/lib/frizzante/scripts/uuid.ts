export function uuid(short = false): string {
    let dt = new Date().getTime()
    let blueprint = "xxxxxxxx-xxxx-yxxx-yxxx-xxxxxxxxxxxx"
    if (short) {
        blueprint = "xyxxyxyx"
    }
    return blueprint.replace(/[xy]/g, function check(c) {
        const r = (dt + Math.random() * 16) % 16 | 0
        dt = Math.floor(dt / 16)
        return (c === "x" ? r : (r & 0x3) | 0x8).toString(16)
    })
}
