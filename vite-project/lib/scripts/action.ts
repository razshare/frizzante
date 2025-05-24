import {getContext} from "svelte";
import type {ServerContext} from "$frizzante/types.ts";
import {route} from "$frizzante/scripts/route.ts";
import {swap} from "./route.ts";

export function action(path: string): {
    method: "POST"
    action: string
    onsubmit: (e: any) => Promise<void>
} {
    const server = getContext("server") as ServerContext<any>
    route(server)
    return {
        method: "POST",
        action: path,
        async onsubmit(e: any) {
            e.preventDefault()
            const form = e.target
            const body = new FormData(form)
            await swap(server,{modifier: "push", method: "GET", path })
            await swap(server,{modifier: "push", method: "POST", path, body})
        }
    }
}