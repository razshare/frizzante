import { IS_BROWSER } from "$lib/scripts/core/is_browser.ts"
import { navigate } from "$lib/scripts/core/navigate.ts"
import { swap } from "$lib/scripts/core/swap"
import { swapping } from "$lib/scripts/core/swapping.ts"
import type { View } from "$lib/scripts/core/view"
import { getContext } from "svelte"
import type { ActionOptions } from "./action_options"
export function action(
    path = "",
    options: ActionOptions = {},
): {
    action: string
    onsubmit: (event: Event) => Promise<void>
} {
    if (!IS_BROWSER) {
        return { action: path, async onsubmit() {} }
    }
    const view: View = getContext("view")
    navigate(view)
    return {
        action: path,
        async onsubmit(event: Event) {
            swapping.active = true
            const pending = setTimeout(function start() {
                if (options.onpending) {
                    options.onpending()
                }
            }, options.pendingDelay ?? 100)
            let error: Error | undefined
            event.preventDefault()
            try {
                const form = event.target as HTMLFormElement
                await swap(form, view).then(function done(record) {
                    record()
                })
            } catch (errorLocal) {
                console.error("swapping failed", errorLocal)
                error = errorLocal as Error
            }
            clearTimeout(pending)
            if (error) {
                if (options.onerror) {
                    options.onerror(error)
                } else {
                    console.error("something went wrong while submitting the form", event.target, error)
                }
            } else {
                if (options.ondone) {
                    options.ondone()
                }
            }
            swapping.active = false
        },
    }
}
