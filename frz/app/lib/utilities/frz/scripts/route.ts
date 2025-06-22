import type { View } from "$lib/utilities/frz/types.ts"
import { IS_BROWSER } from "$lib/utilities/frz/constants.ts"
import { swaps } from "$lib/utilities/frz/scripts/swaps.ts"

let started = false

export function route(view: View<never>): void {
    if (!IS_BROWSER || started) {
        return
    }

    const listener = async function pop(e: PopStateEvent) {
        e.preventDefault()

        const id = e.state ?? ""
        const current = swaps.find(id)

        if (!current) {
            await swaps.swap(view).withPath("/").play(false)
            return
        }

        if (current.position() + 1 != swaps.position()) {
            swaps.teleport(current.position() + 1)
            await current.play(false)
        } else {
            await current.play(true)
        }
    }
    window.addEventListener("popstate", listener)
    started = true
}
