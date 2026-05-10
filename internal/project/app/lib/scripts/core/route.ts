import type { HistoryEntry } from "$lib/scripts/core/history_entry"
import { IS_BROWSER } from "$lib/scripts/core/is_browser.ts"
import { swap } from "$lib/scripts/core/swap"
import type { View } from "$lib/scripts/core/view"
let started = false
export function route(view: View<never>): void {
    if (!IS_BROWSER || started) {
        return
    }
    const form = document.createElement("form")
    const anchor = document.createElement("a")
    const listener = async function pop(e: PopStateEvent) {
        const serialized = (e.state ?? "") as string
        if (serialized !== "") {
            e.preventDefault()
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
        anchor.href = `${window.location}`
        await swap(anchor, view)
    }
    window.addEventListener("popstate", listener)
    started = true
}
