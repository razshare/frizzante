import { getContext } from "svelte"
import type { View } from "$frizzante/core/types.ts"
import { route } from "$frizzante/core/scripts/route.ts"
import { swap } from "$frizzante/core/scripts/swap.ts"
import { IS_BROWSER } from "$frizzante/core/constants.ts"

export function action(path = ""): {
    action: string
    onsubmit: (event: Event) => Promise<void>
} {
    if (!IS_BROWSER) {
        return { action: path, async onsubmit() {} }
    }

    const view = getContext("view") as View<never>
    route(view)
    return {
        action: path,
        async onsubmit(event: Event) {
            event.preventDefault()
            const form = event.target as HTMLFormElement
            await swap(form, view).then(function done(record) {
                record()
                form.reset()
            })
        },
    }
}

