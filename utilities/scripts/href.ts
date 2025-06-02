import {getContext} from "svelte";
import type {View} from "../types.ts";
import {route} from "./route.ts";
import {swaps} from "./swaps.ts";

export function href(path = "", options = {merge: false}): {
    href: string,
    onclick: (e: MouseEvent) => void
} {
    const view = getContext("view") as View<any>
    route(view)
    return {
        href: path,
        async onclick(e: MouseEvent) {
            e.preventDefault()
            await swaps
                .swap(view)
                .withMerge(options.merge)
                .withPath(path)
                .play(true)
            return false
        }
    }
}