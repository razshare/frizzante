import {getContext} from "svelte";
import type {ServerContext} from "$frizzante/types.ts";
import {route} from "$frizzante/scripts/route.ts";
import {swap} from "./route.ts";

export function href(path: string): {
    href: string,
    onclick: (e: MouseEvent) => void
} {
    const server = getContext("server") as ServerContext<any>
    route(server)
    return {
        href: path,
        async onclick(e: MouseEvent) {
            e.preventDefault()
            await swap(server,{modifier: "push", method: "GET", path })
            return false
        }
    }
}