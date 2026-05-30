import type { HrefOptions } from "$lib/scripts/core/href_options"
import { IS_BROWSER } from "$lib/scripts/core/is_browser.ts"
import { navigate } from "$lib/scripts/core/navigate.ts"
import { swap } from "$lib/scripts/core/swap"
import { swapping } from "$lib/scripts/core/swapping.ts"
import type { View } from "$lib/scripts/core/view"
import { getContext } from "svelte"
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
    const view: View = getContext("view")
    navigate(view)
    return {
        href: path,
        async onclick(event: MouseEvent) {
            swapping.active = true
            let error: Error | undefined
            const pending = setTimeout(function start() {
                if (options.onpending) {
                    options.onpending()
                }
            }, options.pendingDelay ?? 100)
            event.preventDefault()
            try {
                const record = await swap(anchor, view)
                record()
            } catch (errorLocal) {
                console.error("swapping failed", errorLocal)
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
            swapping.active = false
            return false
        },
    }
}
