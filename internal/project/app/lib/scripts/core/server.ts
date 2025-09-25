import { render as _render } from "svelte/server"
import Router from "$lib/components/core/Router.svelte"
export async function render(args: Record<string, never>) {
    // @ts-expect-error
    return _render(Router, { props: args })
}
