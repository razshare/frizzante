import {getContext} from "svelte";
import {navigate} from "$frizzante/scripts/route.ts";
import type {ServerContext} from "$frizzante/types.ts";

export function action(id: string): {
    method: "POST"
    action: string
    onsubmit: (e: any) => Promise<void>
} {
    let server = getContext("server") as ServerContext<any>
    return {
        method: "POST",
        action: server.ids[id],
        async onsubmit(e: any) {
            e.preventDefault()
            const form = e.target
            const body = new FormData(form)
            const method = form.method.toUpperCase()
            const headers = {"Accept": "application/json"}
            const response = await fetch(form.action, {method, headers, body})
            if (response.status >= 300) {
                return
            }

            const json = await response.json()

            server.data = json.data
            server.ids = json.ids

            if (server.id !== json.id) {
                navigate(server, json.id, server.data)
                    .then(function done() {
                        server.id = json.id
                    })
            }
        }
    }
}