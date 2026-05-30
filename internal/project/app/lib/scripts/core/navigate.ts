import type { HistoryEntry } from "$lib/scripts/core/history_entry"
import { IS_BROWSER } from "$lib/scripts/core/is_browser.ts"
import { swap } from "$lib/scripts/core/swap"
import type { View } from "$lib/scripts/core/view"
let started = false
export function navigate(view: View): void {
    if (!IS_BROWSER || started) {
        return
    }
    const form = document.createElement("form")
    const anchor = document.createElement("a")
    const listener = async function pop(event: PopStateEvent) {
        const serialized: string = event.state ?? ""
        if (serialized !== "") {
            event.preventDefault()
            const entry: HistoryEntry = JSON.parse(serialized)
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
