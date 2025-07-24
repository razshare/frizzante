import { getContext } from "svelte"
import type { View } from "$frizzante/core/types.ts"
import { route } from "$frizzante/core/scripts/route.ts"
import { swaps } from "$frizzante/core/scripts/swaps.ts"

export function action(path = ""): {
    action: string
    onsubmit: (event: Event) => Promise<void>
} {
    const view = getContext("view") as View<never>
    route(view)
    return {
        action: path,
        async onsubmit(event: Event) {
            event.preventDefault()
            const form = event.target as HTMLFormElement
            const body = new FormData(form)
            const target = event.target as HTMLFormElement

            await swaps
                .swap(view)
                .withMethod(target.method)
                .withPath(path)
                .withBody(body)
                .withUpdate(true)
                .play()
                .then(function done() {
                    form.reset()
                })
        },
    }
}
