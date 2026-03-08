export function delay(milliseconds: number): Promise<void> {
    return new Promise(function run(stop) {
        setTimeout(stop, milliseconds)
    })
}
