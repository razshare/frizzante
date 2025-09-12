import { getContext } from "svelte"
import type { View } from "$lib/scripts/core/types.ts"
import { route } from "$lib/scripts/core/route.ts"
import { swap } from "$lib/scripts/core/swap.ts"
import { IS_BROWSER } from "$lib/scripts/core/constants.ts"

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
