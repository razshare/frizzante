import { mount } from "svelte"
import Router from "$lib/components/core/ClientRouter.svelte"
export function render(target: HTMLElement, args: Record<string, never>) {
    target.innerHTML = ""
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    mount(Router, { target, props: args })
}
