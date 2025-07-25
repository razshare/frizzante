import type { View } from "$frizzante/core/types.ts"
import { uuid } from "$frizzante/core/scripts/uuid.ts"

type SwapAction = {
    method: () => string
    path: () => string
    body: () => unknown
    position: () => number
    withMethod: (method: string) => SwapAction
    withPath: (path: string) => SwapAction
    withBody: (body: FormData) => SwapAction
    withUpdate: (update: boolean) => SwapAction
    play: () => Promise<void>
}

let nextPosition = 0
const record = {} as Record<string, SwapAction>

function find(id: string): false | SwapAction {
    return record[id] ?? false
}

function swap(view: View<unknown>): SwapAction {
    const swapPosition = nextPosition++
    let swapMethod = "GET"
    let swapPath = location.pathname
    let swapUpdate = false
    let swapBody: FormData

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
        withUpdate(update: boolean): SwapAction {
            swapUpdate = update
            return this
        },
        async play() {
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
            view.renderMode = json.renderMode

            if (swapUpdate) {
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
