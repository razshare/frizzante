import type { HistoryEntry, View } from "$frizzante/core/types.ts"

let lastView: false | string = false

export async function swap(
    target: HTMLAnchorElement | HTMLFormElement,
    view: View<unknown>,
): Promise<() => void> {
    if (lastView === false) {
        lastView = location.toString()
    }

    let res: Response
    let method: "GET" | "POST" = "GET"
    const body: Record<string, string> = {}

    if (target.nodeName === "A") {
        const anchor = target as HTMLAnchorElement
        res = await fetch(anchor.href, {
            headers: {
                Accept: "application/json",
            },
        })
    } else if (target.nodeName === "FORM") {
        const form = target as HTMLFormElement
        const data = new FormData(form)
        const params = new URLSearchParams()
        let query = ""

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
            res = await fetch(`${form.action}${query}`, {
                headers: {
                    Accept: "application/json",
                },
            })
        } else {
            res = await fetch(form.action, {
                method,
                body: data,
                headers: {
                    Accept: "application/json",
                },
            })
        }
    } else {
        return function push() {}
    }

    const txt = await res.text()

    if ("" === txt) {
        return function push() {}
    }

    const json = JSON.parse(txt)

    view.data = json.data
    view.name = json.name
    view.renderMode = json.renderMode

    const sameView = lastView === res.url
    lastView = res.url

    return function push() {
        if (sameView) {
            return
        }

        const entry: HistoryEntry = {
            nodeName: target.nodeName,
            method,
            url: res.url,
            body,
        }

        window.history.pushState(JSON.stringify(entry), "", res.url)
    }
}
