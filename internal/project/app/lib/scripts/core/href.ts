import { getContext } from "svelte"
import type { View } from "$lib/scripts/core/view"
import { route } from "$lib/scripts/core/route.ts"
import { swap } from "$lib/scripts/core/swap.ts"
import { IS_BROWSER } from "$lib/scripts/core/is_browser.ts"
import { swapping } from "$lib/scripts/core/swapping.ts"
export function href(path = ""): {
    href: string
    onclick: (event: MouseEvent) => Promise<boolean>
} {
    if (!IS_BROWSER) {
        return {
            href: path,
            async onclick() {
                return true
            },
        }
    }
    const anchor = document.createElement("a")
    anchor.href = path
    const view = getContext("view") as View<never>
    route(view)
    return {
        href: path,
        async onclick(event: MouseEvent) {
            swapping.active = true
            event.preventDefault()
            try {
                const record = await swap(anchor, view)
                record()
            } catch (error) {
                console.error("swapping failed", error)
            }
            swapping.active = false
            return false
        },
    }
}
