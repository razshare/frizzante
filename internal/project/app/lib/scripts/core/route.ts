import type { HistoryEntry, View } from "$lib/scripts/core/types.ts"
import { IS_BROWSER } from "$lib/scripts/core/constants.ts"
import { swap } from "$lib/scripts/core/swap.ts"

let started = false

export function route(view: View<never>): void {
    if (!IS_BROWSER || started) {
        return
    }

    const form = document.createElement("form")
    const anchor = document.createElement("a")

    const listener = async function pop(e: PopStateEvent) {
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
