import {getContext} from "svelte";
import type {View} from "../types.ts";
import {route} from "./route.ts";
import {swaps} from "./swaps.ts";

export function action(path = "", options = {merge: false}): {
    method: "POST"
    action: string
    onsubmit: (e: any) => Promise<void>
} {
    const view = getContext("view") as View<any>
    route(view)
    return {
        method: "POST",
        action: path,
        async onsubmit(e: any) {
            e.preventDefault()
            const form = e.target as HTMLFormElement
            const body = new FormData(form)

            await swaps
                .swap(view)
                .withMerge(options.merge)
                .withMethod("POST")
                .withPath(path)
                .withBody(body)
                .play(true).then(function done(){
                    form.reset()
                })
        }
    }
}