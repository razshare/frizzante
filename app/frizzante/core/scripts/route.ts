import type { View } from "$frizzante/core/types.ts"
import { IS_BROWSER } from "$frizzante/core/constants.ts"
import { find, swap } from "$frizzante/core/scripts/swap.ts"

let started = false

export function route(view: View<never>): void {
    if (!IS_BROWSER || started) {
        return
    }

    const base = location.pathname

    const listener = async function pop(e: PopStateEvent) {
        e.preventDefault()

        const id = (e.state ?? -1) as number

        if (id >= 0) {
            const config = find(id)
            if (!config) {
                console.warn("swap configuration not found", { id })
                return
            }
            await swap(config)
        } else {
            await swap({ method: "GET", path: base, view, body: false })
        }
    }
    window.addEventListener("popstate", listener)
    started = true
}
