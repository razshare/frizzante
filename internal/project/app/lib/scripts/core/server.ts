import ServerRouter from "$lib/components/core/server_router.svelte"
import { render as ssr } from "svelte/server"
export async function render(props: any) {
    return ssr(ServerRouter, { props })
}
