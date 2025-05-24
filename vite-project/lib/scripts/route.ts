import type {ServerContext} from "../types.ts";
import {swaps} from "./swaps.ts";

let started = false
const IS_BROWSER = typeof document !== 'undefined'

export function route(server: ServerContext<any>): void {
    if (!IS_BROWSER || started) {
        return
    }

    const listener = async function pop(e: PopStateEvent) {
        e.preventDefault();

        const id = e.state ?? ""
        const current = swaps.find(id)

        if (!current) {
            await swaps.swap(server).withPath("/").play(false)
            return
        }

        if (current.position() + 1 != swaps.position()) {
            swaps.teleport(current.position() + 1)
            await current.play(false)
        } else {
            await current.play(true)
        }
    }
    window.addEventListener("popstate", listener);
    started = true
}
