import {getContext} from "svelte";
import type {ServerContext} from "../types.ts";
import {route} from "./route.ts";
import {swaps} from "./swaps.ts";

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

            await swaps
                .swap(server)
                .withMethod("POST")
                .withPath(path)
                .withBody(body)
                .play(true)
        }
    }
}