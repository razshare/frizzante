import { render as _render } from "svelte/server"
import ServerRouter from "$lib/components/core/server_router.svelte"
export async function render(args: Record<string, never>) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    return _render(ServerRouter, { props: args })
}
