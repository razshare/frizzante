import type { HistoryEntry, View } from "$lib/scripts/core/types"

let lastUrl: false | string = false

export async function swap(target: HTMLAnchorElement | HTMLFormElement, view: View<unknown>): Promise<() => void> {
    if (lastUrl === false) {
        lastUrl = location.toString()
    }

    let response: Response
    let method: "GET" | "POST" = "GET"
    const body: Record<string, string> = {}

    if (target.nodeName === "A") {
        const anchor = target as HTMLAnchorElement
        response = await fetch(anchor.href, {
            headers: {
                Accept: "application/json",
            },
        })
    } else if (target.nodeName === "FORM") {
        const form = target as HTMLFormElement
        const data = new FormData(form)
        const params = new URLSearchParams()
        let query = ""

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
            query = `${params.toString()}`
            if (query !== "") {
                if (form.action.includes("?")) {
                    query = "&" + query
                } else {
                    query = "?" + query
                }
            }
            response = await fetch(`${form.action}${query}`, {
                headers: {
                    Accept: "application/json",
                },
            })
        } else {
            response = await fetch(form.action, {
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

    view.pending = true
    const text = await response.text()

    if ("" === text) {
        return function push() {}
    }

    const remote = JSON.parse(text) as View<Record<string, unknown>>

    await view.snapshot()
    view.name = remote.name
    view.align = remote.align
    view.render = remote.render
    if (view.align === 1) {
        if (typeof view.props != "object") {
            console.warn("view alignment intends to merge props, but local view props is not an object")
            // Noop.
        } else if (typeof remote.props != "object") {
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
    view.pending = false

    const stationary = lastUrl === response.url
    lastUrl = response.url

    return function push() {
        if (stationary) {
            return
        }

        const entry: HistoryEntry = {
            nodeName: target.nodeName,
            method,
            url: response.url,
            body,
        }

        window.history.pushState(JSON.stringify(entry), "", response.url)
    }
}
