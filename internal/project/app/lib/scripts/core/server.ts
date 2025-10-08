import { render as _render } from "svelte/server"
import Router from "$lib/components/core/ServerRouter.svelte"
export async function render(args: Record<string, never>) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    return _render(Router, { props: args })
}
