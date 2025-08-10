import type { HistoryEntry, View } from "$frizzante/core/types.ts"
import { IS_BROWSER } from "$frizzante/core/constants.ts"
import { find, swap } from "$frizzante/core/scripts/swap.ts"

let started = false

export function route(view: View<never>): void {
    if (!IS_BROWSER || started) {
        return
    }

    const listener = async function pop(e: PopStateEvent) {
        e.preventDefault()

        const entry = (e.state ?? false) as false | HistoryEntry

        if (entry) {
            let config = find(entry.id)
            if (!config) {
                console.info(
                    "swap configuration not found, reconstructing it",
                    { id: entry.id },
                )
                config = {
                    method: entry.method,
                    path: entry.path,
                    body: false,
                    view,
                }
            }
            await swap(config)
        } else {
            await swap({ method: "GET", path: "/", view, body: false })
        }
    }
    window.addEventListener("popstate", listener)
    started = true
}
