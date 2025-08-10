import type { View } from "$frizzante/core/types.ts"
import { IS_BROWSER } from "$frizzante/core/constants.ts"
import { swap } from "$frizzante/core/scripts/swap.ts"

let started = false

export function route(view: View<never>): void {
    if (!IS_BROWSER || started) {
        return
    }

    const anchor = document.createElement("a")

    const listener = async function pop(e: PopStateEvent) {
        e.preventDefault()
        anchor.href = (e.state ?? "/") as string
        await swap(anchor, view)
    }

    window.addEventListener("popstate", listener)
    started = true
}
