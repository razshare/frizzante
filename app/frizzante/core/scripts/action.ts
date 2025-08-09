import { getContext } from "svelte"
import type { View } from "$frizzante/core/types.ts"
import { route } from "$frizzante/core/scripts/route.ts"
import { swap } from "$frizzante/core/scripts/swap.ts"

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
            await swap({ method: target.method, path, body, view }).then(
                function done(record) {
                    record()
                    form.reset()
                },
            )
        },
    }
}
