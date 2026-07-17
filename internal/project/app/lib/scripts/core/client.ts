import ClientRouter from "$lib/components/core/client_router.svelte"
import type { ServerRouterProps } from "$lib/types/core/server_router_props"
import { hydrate } from "svelte"
export function render(target: HTMLElement, props: ServerRouterProps) {
    target.innerText = ""
    hydrate(ClientRouter, { target, props })
}
