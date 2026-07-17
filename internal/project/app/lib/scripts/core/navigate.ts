import { IS_BROWSER } from "$lib/scripts/core/is_browser.ts"
import { swap } from "$lib/scripts/core/swap.svelte"
import type { HistoryEntry } from "$lib/types/core/history_entry"
let started = false
export function navigate(): void {
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
                await swap(anchor)
            }
            form.innerHTML = ""
            for (const key in entry.body) {
                const value = entry.body[key]
                const input = document.createElement("input")
                input.value = value
                form.appendChild(input)
            }
            await swap(form)
            return
        }
        anchor.href = `${window.location}`
        await swap(anchor)
    }
    window.addEventListener("popstate", listener)
    started = true
}
