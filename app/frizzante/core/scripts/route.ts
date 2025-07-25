import type { View } from "$frizzante/core/types.ts"
import { IS_BROWSER } from "$frizzante/core/constants.ts"
import { swaps } from "$frizzante/core/scripts/swaps.ts"

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
            await swaps.swap(view).withPath("/").play()
            return
        }

        if (current.position() + 1 != swaps.position()) {
            swaps.teleport(current.position() + 1)
            await current.play()
        } else {
            await current.withUpdate(true).play()
        }
    }
    window.addEventListener("popstate", listener)
    started = true
}
