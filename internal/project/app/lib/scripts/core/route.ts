import type { View } from "$lib/scripts/core/view"
import type { HistoryEntry } from "$lib/scripts/core/history_entry"
import { IS_BROWSER } from "$lib/scripts/core/is_browser.ts"
import { swap } from "$lib/scripts/core/swap.ts"
import { swapping } from "./swapping"
let started = false
export function route(view: View<never>): void {
    if (!IS_BROWSER || started) {
        return
    }
    const form = document.createElement("form")
    const anchor = document.createElement("a")
    const listener = async function pop(e: PopStateEvent) {
        // we don't want to interfere with popstate events
        // that are not triggered by href() and action().
        if (!swapping.active) {
            return
        }
        e.preventDefault()
        const serialized = (e.state ?? "") as string
        if (serialized !== "") {
            const entry = JSON.parse(serialized) as HistoryEntry
            if (entry.method === "GET") {
                anchor.href = entry.url
                await swap(anchor, view)
            }
            form.innerHTML = ""
            for (const key in entry.body) {
                const value = entry.body[key]
                const input = document.createElement("input")
                input.value = value
                form.appendChild(input)
            }
            await swap(form, view)
            return
        }
        anchor.href = "/"
        await swap(anchor, view)
    }
    window.addEventListener("popstate", listener)
    started = true
}
