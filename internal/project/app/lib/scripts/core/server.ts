import { render as ssr } from "svelte/server"
import Router from "$lib/components/core/Router.svelte"
export async function render(args: Record<string, never>) {
    return ssr(Router, { props: args })
}
