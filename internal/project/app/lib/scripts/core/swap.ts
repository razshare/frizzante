import type { View } from "$lib/scripts/core/view"
import type { HistoryEntry } from "$lib/scripts/core/history_entry"
let lastUrl: false | string = false
export async function swap(target: HTMLAnchorElement | HTMLFormElement, view: View<unknown>): Promise<() => void> {
    if (lastUrl === false) {
        lastUrl = location.toString()
    }
    let requestUrl: string
    let response: Response
    let method: "GET" | "POST" = "GET"
    const body: Record<string, string> = {}
    if (target.nodeName === "A") {
        const anchor = target as HTMLAnchorElement
        requestUrl = anchor.href
        console.log({ requestUrl })
        if (view.type === "snapshot") {
            requestUrl = requestUrl.replace(/\/+$/, "") + "/data.json"
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
        data.forEach(function each(value, key) {
            if (value instanceof File) {
                return
            }
            body[key] = `${value}`
            params.append(key, `${value}`)
        })
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
    const remote = JSON.parse(text) as View<Record<string, unknown>>
    await view.pin()
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
    const hash = requestUrl.split("#", 2)[1] ?? ""
    const stationary = lastUrl === fixedResponseUrl
    lastUrl = fixedResponseUrl
    return function push() {
        if (stationary) {
            return
        }
        const entry: HistoryEntry = {
            nodeName: target.nodeName,
            method,
            url: fixedResponseUrl,
            body,
        }
        console.log({ requestUrl, hash })
        if (hash !== "") {
            window.history.pushState(JSON.stringify(entry), "", `${fixedResponseUrl}#${hash}`)
        } else {
            window.history.pushState(JSON.stringify(entry), "", fixedResponseUrl)
        }
    }
}
