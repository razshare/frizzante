import { getContext } from "svelte"
import type { View } from "$lib/utilities/frizzante/types.ts"
import { route } from "$lib/utilities/frizzante/scripts/route.ts"
import { swaps } from "$lib/utilities/frizzante/scripts/swaps.ts"

export function href(path = ""): {
    href: string
    onclick: (e: MouseEvent) => void
} {
    const view = getContext("view") as View<never>
    route(view)
    return {
        href: path,
        async onclick(e: MouseEvent) {
            e.preventDefault()
            await swaps.swap(view).withPath(path).play(true)
            return false
        },
    }
}
