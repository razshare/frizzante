import { getContext } from "svelte"
import type { View } from "$lib/utilities/types.ts"
import { route } from "$lib/utilities/scripts/route.ts"
import { swaps } from "$lib/utilities/scripts/swaps.ts"

export function action(path = ""): {
    action: string
    onsubmit: (e: Event) => Promise<void>
} {
    const view = getContext("view") as View<never>
    route(view)
    return {
        action: path,
        async onsubmit(e: Event) {
            e.preventDefault()
            const form = e.target as HTMLFormElement
            const body = new FormData(form)
            const target = e.target as HTMLFormElement

            await swaps
                .swap(view)
                .withMethod(target.method)
                .withPath(path)
                .withBody(body)
                .play(true)
                .then(function done() {
                    form.reset()
                })
        },
    }
}
