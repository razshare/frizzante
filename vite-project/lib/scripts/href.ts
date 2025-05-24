import {getContext} from "svelte";
import type {ServerContext} from "../types.ts";
import {route} from "./route.ts";
import {swaps} from "./swaps.ts";

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
            await swaps.swap(server).withPath(path).play(true)
            return false
        }
    }
}