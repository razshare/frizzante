import ServerRouter from "$lib/components/core/server_router.svelte"
import type { ClientRouterProps } from "$lib/types/core/client_router_props"
import { render as ssr } from "svelte/server"
export async function render(props: ClientRouterProps) {
    return ssr(ServerRouter, { props })
}
