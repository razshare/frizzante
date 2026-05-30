import ClientRouter from "$lib/components/core/client_router.svelte"
import type { View } from "$lib/scripts/core/view.d.ts"
import { mount } from "svelte"
export function render(target: HTMLElement, args: View) {
    target.innerText = ""
    mount(ClientRouter, { target, props: args })
}
