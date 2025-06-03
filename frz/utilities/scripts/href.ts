import {getContext} from "svelte";
import type {View} from "../types.ts";
import {route} from "./route.ts";
import {swaps} from "./swaps.ts";

export function href(path = ""): {
    href: string,
    onclick: (e: MouseEvent) => void
} {
    const view = getContext("view") as View<never>
    route(view)
    return {
        href: path,
        async onclick(e: MouseEvent) {
            e.preventDefault()
            await swaps
                .swap(view)
                .withPath(path)
                .play(true)
            return false
        }
    }
}