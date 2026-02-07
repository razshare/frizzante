import { mount } from "svelte"
import ClientRouter from "$lib/components/core/client_router.svelte"
import type { View } from "$lib/scripts/core/view.d.ts"
export function render(target: HTMLElement, args: View<Record<string, unknown>>) {
    target.innerHTML = ""
    mount(ClientRouter, { target, props: args })
}
