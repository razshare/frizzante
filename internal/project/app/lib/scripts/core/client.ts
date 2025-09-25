import { mount  } from "svelte"
import Router from "$lib/components/core/Router.svelte"
export function render(target: HTMLElement, args: Record<string, never>) {
    // @ts-expect-error
    mount(Router, { target, props: args })
}
