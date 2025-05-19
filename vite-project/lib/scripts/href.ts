import {getContext} from "svelte";
import {navigate} from "$frizzante/scripts/route.ts";
import type {ServerContext} from "$frizzante/types.ts";

export function href(to: string): {
    href: string,
    onclick: (e: MouseEvent) => void
} {
    const server = getContext("server") as ServerContext<any>
    return {
        href: server.ids[to],
        async onclick(e: MouseEvent) {
            e.preventDefault()
            await navigate(server, to)
            return false
        }
    }
}