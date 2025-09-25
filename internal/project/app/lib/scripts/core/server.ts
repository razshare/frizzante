import { render as _render } from "svelte/server"
import Router from "$lib/components/core/ServerRouter.svelte"
export async function render(args: Record<string, never>) {
    return _render(Router, { props: args })
}
