import { mount } from "svelte"
import ClientRouter from "$lib/components/core/client_router.svelte"
export function render(target: HTMLElement, args: Record<string, never>) {
    target.innerHTML = ""
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    mount(ClientRouter, { target, props: args })
}
