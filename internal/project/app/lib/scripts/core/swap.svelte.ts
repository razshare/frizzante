import type { views as clientViews } from "$exports.client"
import type { views as serverViews } from "$exports.server"
import { root } from "$lib/scripts/core/root.svelte"
import type { HistoryEntry } from "$lib/types/core/history_entry"
import { SvelteURLSearchParams } from "svelte/reactivity"
let lastUrl: false | string = false
export async function swap(target: HTMLAnchorElement | HTMLFormElement): Promise<() => void> {
    if (lastUrl === false) {
        lastUrl = location.toString()
    }
    let requestUrl: string
    let response: Response
    let method: "GET" | "POST" = "GET"
    let hash = ""
    const body: Record<string, string> = {}
    if (target.nodeName === "A") {
        const anchor = target as HTMLAnchorElement
        const parts = anchor.href.split("#", 2)
        requestUrl = parts[0]
        if (root.type === "snapshot") {
            requestUrl = requestUrl.replace(/\/+$/, "") + "/data.json"
        }
        if (parts.length >= 2) {
            hash = `#${parts[1]}`
        }
        response = await fetch(requestUrl, {
            headers: {
                Accept: "application/json",
            },
        })
    } else if (target.nodeName === "FORM") {
        const form = target as HTMLFormElement
        const data = new FormData(form)
        const params = new SvelteURLSearchParams()
        requestUrl = form.action.split("?")[0] ?? ""
        if (root.type === "snapshot") {
            requestUrl = requestUrl.replace(/\/+$/, "") + "/data.json"
        }
        form.reset()
        for (const key of data.keys()) {
            const value = data.get(key)
            if (value instanceof File) {
                continue
            }
            body[key] = `${value}`
            params.append(key, `${value}`)
        }

        method = form.method.toUpperCase() as "GET" | "POST"
        if (method === "GET") {
            const query = `?${params.toString()}`
            response = await fetch(`${requestUrl}${query}`, {
                headers: {
                    Accept: "application/json",
                },
            })
        } else {
            requestUrl = form.action
            if (root.type === "snapshot") {
                requestUrl += "/data.json"
            }
            response = await fetch(requestUrl, {
                method,
                body: data as unknown as BodyInit,
                headers: {
                    Accept: "application/json",
                },
            })
        }
    } else {
        return function push(): void {}
    }
    const text = await response.text()
    if (text === "") {
        return function push(): void {}
    }
    const remote = JSON.parse(text) as {
        name: keyof typeof clientViews | keyof typeof serverViews
        render: number
        props: Record<string, unknown>
    }
    root.view = {
        name: remote.name,
        ondone() {
            return remote.props
        },
    }
    let fixedResponseUrl = response.url
    if (root.type === "snapshot") {
        fixedResponseUrl = fixedResponseUrl.replace(/\/data\.json$/, "")
    }
    lastUrl = fixedResponseUrl
    return function push() {
        const entry: HistoryEntry = {
            nodeName: target.nodeName,
            method,
            url: fixedResponseUrl,
            body,
        }
        window.history.pushState(JSON.stringify(entry), "", `${fixedResponseUrl}${hash}`)
    }
}
