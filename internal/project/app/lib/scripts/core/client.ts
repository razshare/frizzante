import ClientRouter from "$lib/components/core/client_router.svelte"
import { hydrate } from "svelte"
export function render(target: HTMLElement, props: any) {
    target.innerText = ""
    hydrate(ClientRouter, { target, props })
}
