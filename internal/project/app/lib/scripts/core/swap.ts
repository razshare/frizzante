import type { HistoryEntry } from "$lib/scripts/core/history_entry"
import type { View } from "$lib/scripts/core/view"
let lastUrl: false | string = false
export async function swap(target: HTMLAnchorElement | HTMLFormElement, view: View): Promise<() => void> {
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
        if (view.type === "snapshot") {
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
        const params = new URLSearchParams()
        requestUrl = form.action.split("?")[0] ?? ""
        if (view.type === "snapshot") {
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
            if (view.type === "snapshot") {
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
    const remote = JSON.parse(text) as View
    view.name = remote.name
    view.align = remote.align
    view.render = remote.render
    if (view.align === 1) {
        if (typeof view.props !== "object") {
            console.warn("view alignment intends to merge props, but local view props is not an object")
            // Noop.
        } else if (typeof remote.props !== "object") {
            console.warn("view alignment intends to merge props, but remote props is not an object")
            // Noop.
        } else {
            view.props = {
                ...view.props,
                ...remote.props,
            }
        }
    } else {
        view.props = remote.props
    }
    let fixedResponseUrl = response.url
    if (view.type === "snapshot") {
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
