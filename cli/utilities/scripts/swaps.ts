import type { View } from "$lib/utilities/types.ts"
import { uuid } from "$lib/utilities/scripts/uuid.ts"

type SwapAction = {
    method: () => string
    path: () => string
    body: () => unknown
    position: () => number
    withMethod: (method: string) => SwapAction
    withPath: (path: string) => SwapAction
    withBody: (body: FormData) => SwapAction
    play: (update: boolean) => Promise<void>
}

let nextPosition = 0
const record = {} as Record<string, SwapAction>

function find(id: string): false | SwapAction {
    return record[id] ?? false
}

function swap(view: View<unknown>): SwapAction {
    let swapMethod = "GET"
    let swapPath = location.pathname
    let swapBody: FormData
    const swapPosition = nextPosition++

    return {
        method() {
            return swapMethod
        },
        path() {
            return swapPath
        },
        body() {
            return swapBody
        },
        position() {
            return swapPosition
        },
        withMethod(method: string) {
            swapMethod = method.toUpperCase()
            return this
        },
        withPath(path: string) {
            swapPath = path
            return this
        },
        withBody(body: FormData) {
            swapBody = body
            return this
        },

        async play(update: boolean) {
            const payload = {
                method: swapMethod,
                headers: { Accept: "application/json" },
            } as RequestInit

            let query = ""

            if ("GET" === swapMethod) {
                if (swapBody && typeof swapBody === "object") {
                    const params = new URLSearchParams()
                    swapBody.forEach(function each(value, key) {
                        params.append(key, `${value}`)
                    })

                    query = `${params.toString()}`

                    if (swapPath.includes("?")) {
                        query = "&" + query
                    } else {
                        query = "?" + query
                    }
                }
            } else {
                payload.body = swapBody as BodyInit
            }

            const response = await fetch(`${swapPath}${query}`, payload)

            const text = await response.text()

            if ("" === text) {
                return
            }

            const json = JSON.parse(text)

            view.data = json.data
            view.name = json.name
            view.error = json.error

            if (update) {
                const id = uuid()
                record[id] = this
                window.history.pushState(id, "", response.url)
            }
        },
    }
}

function position(): number {
    return nextPosition
}

function teleport(position: number) {
    nextPosition = position
}

export const swaps = {
    swap,
    find,
    position,
    teleport,
}
