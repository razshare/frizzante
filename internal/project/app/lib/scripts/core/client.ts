import { mount as csr } from "svelte"
import Router from "$lib/components/core/Router.svelte"
export function render(target: HTMLElement, args: Record<string, never>) {
    csr(Router, { target, props: args })
}
