import { getContext } from "svelte"
import type { View } from "$frizzante/core/types.ts"
import { route } from "$frizzante/core/scripts/route.ts"
import { swap } from "$frizzante/core/scripts/swap.ts"
import { IS_BROWSER } from "$frizzante/core/constants.ts"

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
            event.preventDefault()
            const record = await swap(anchor, view)
            record()
            return false
        },
    }
}

