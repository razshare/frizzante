import ServerRouter from "$lib/components/core/server_router.svelte"
import type { View } from "$lib/scripts/core/view.d.ts"
import { render as _render } from "svelte/server"
export async function render(args: View) {
    return _render(ServerRouter, { props: args })
}
