import { IS_BROWSER } from "$lib/scripts/core/is_browser.ts"
import { root } from "$lib/scripts/core/root.svelte"
import { swap } from "$lib/scripts/core/swap.svelte"
import type { HrefOptions } from "$lib/types/core/href_options"
export function href(
    path = "",
    options: HrefOptions = {},
): {
    href: string
    onclick: (event: MouseEvent) => Promise<boolean>
} {
    if (!IS_BROWSER) {
        return {
            href: path,
            async onclick() {
                return true
            },
        }
    }
    const anchor = document.createElement("a")
    anchor.href = path
    return {
        href: path,
        async onclick(event: MouseEvent) {
            root.swapping = true
            let error: Error | undefined
            const pending = setTimeout(function start() {
                if (options.onpending) {
                    options.onpending()
                }
            }, options.pendingDelay ?? 100)
            event.preventDefault()
            try {
                const record = await swap(anchor)
                record()
            } catch (errorLocal) {
                error = errorLocal as Error
            }
            clearTimeout(pending)
            if (error) {
                if (options.onerror) {
                    options.onerror(error)
                }
            } else {
                if (options.ondone) {
                    options.ondone()
                }
            }
            root.swapping = false
            return false
        },
    }
}
