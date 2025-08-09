import { getContext } from "svelte"
import type { View } from "$frizzante/core/types.ts"
import { route } from "$frizzante/core/scripts/route.ts"
import { swap } from "$frizzante/core/scripts/swap.ts"

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
            const record = await swap({
                method: "GET",
                path,
                body: false,
                view,
            })
            record()
            return false
        },
    }
}
