import type { HistoryEntry, SwapConfig } from "$frizzante/core/types.ts"

let id = 0
const configs: Record<number, SwapConfig> = {}

export function find(id: number): false | SwapConfig {
    return configs[id] ?? false
}

export async function swap(config: SwapConfig): Promise<() => void> {
    const payload = {
        method: config.method.toUpperCase(),
        headers: { Accept: "application/json" },
    } as RequestInit

    let query = ""
    let pushState = true

    if ("GET" === config.method.toUpperCase()) {
        if (config.body && typeof config.body === "object") {
            const params = new URLSearchParams()
            config.body.forEach(function each(value, key) {
                params.append(key, `${value}`)
            })

            query = `${params.toString()}`

            if (config.path.includes("?")) {
                query = "&" + query
            } else {
                query = "?" + query
            }
        }
    } else if (config.body) {
        payload.body = config.body as BodyInit
        pushState = false
    }

    const res = await fetch(`${config.path}${query}`, payload)

    if (res.redirected) {
        pushState = false
    }

    const txt = await res.text()

    if ("" === txt) {
        return function push() {}
    }

    const json = JSON.parse(txt)

    config.view.data = json.data
    config.view.name = json.name
    config.view.renderMode = json.renderMode
    if (pushState) {
        configs[++id] = config
    }

    return function push() {
        if (!pushState) {
            return
        }

        const friendly: HistoryEntry = {
            id,
            method: config.method,
            path: config.path,
        }

        window.history.pushState(friendly, "", res.url)
    }
}
