import { IS_BROWSER } from "$lib/scripts/core/is_browser.ts"
import { root } from "$lib/scripts/core/root.svelte"
import { swap } from "$lib/scripts/core/swap.svelte"
import type { ActionOptions } from "$lib/types/core/action_options"
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
    return {
        action: path,
        async onsubmit(event: Event) {
            root.swapping = true
            const pending = setTimeout(function start() {
                if (options.onpending) {
                    options.onpending()
                }
            }, options.pendingDelay ?? 100)
            let error: Error | undefined
            event.preventDefault()
            try {
                const form = event.target as HTMLFormElement
                await swap(form).then(function done(record) {
                    record()
                })
            } catch (errorLocal) {
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
            root.swapping = false
        },
    }
}
