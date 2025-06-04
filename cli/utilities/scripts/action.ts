import {getContext} from "svelte";
import type {View} from "../types.ts";
import {route} from "./route.ts";
import {swaps} from "./swaps.ts";

export function action(path = ""): {
    method: "POST"
    action: string
    onsubmit: (e: Event) => Promise<void>
} {
    const view = getContext("view") as View<never>
    route(view)
    return {
        method: "POST",
        action: path,
        async onsubmit(e: Event) {
            e.preventDefault()
            const form = e.target as HTMLFormElement
            const body = new FormData(form)

            await swaps
                .swap(view)
                .withMethod("POST")
                .withPath(path)
                .withBody(body)
                .play(true).then(function done() {
                    form.reset()
                })
        }
    }
}