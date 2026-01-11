import type { HistoryEntry, View } from "$lib/scripts/core/types"

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
        if(view.isSnapshot) {
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
        let query = ""

        requestUrl = form.action.split("?")[0] ?? ""
        if(view.isSnapshot) {
            requestUrl = requestUrl.replace(/\/+$/,"") + "/data.json"
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
            query = `?${params.toString()}`
            response = await fetch(`${requestUrl}${query}`, {
                headers: {
                    Accept: "application/json",
                },
            })
        } else {
            requestUrl = form.action
            if(view.isSnapshot) {
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
        return function push() {}
    }

    const text = await response.text()

    if (text === "") {
        return function push() {}
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

    if(view.isSnapshot) {
        fixedResponseUrl = fixedResponseUrl.replace(/\/data\.json$/, "")
    }

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

        window.history.pushState(JSON.stringify(entry), "", fixedResponseUrl)
    }
}
