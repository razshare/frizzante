import { getContext } from "svelte"
import type { View } from "$frizzante/core/types.ts"
import { route } from "$frizzante/core/scripts/route.ts"
import { swaps } from "$frizzante/core/scripts/swaps.ts"

export function href(path = ""): {
    href: string
    onclick: (event: MouseEvent) => Promise<boolean>
} {
    const view = getContext("view") as View<never>
    route(view)
    return {
        href: path,
        async onclick(event: MouseEvent) {
            event.preventDefault()
            await swaps.swap(view).withPath(path).withUpdate(true).play()
            return false
        },
    }
}
