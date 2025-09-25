import { mount } from "svelte"
import Router from "$lib/components/core/ClientRouter.svelte"
export function render(target: HTMLElement, args: Record<string, never>) {
    target.innerHTML = ""
    mount(Router, { target, props: args })
}
