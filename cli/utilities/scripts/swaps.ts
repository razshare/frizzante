import type {View} from "../types.ts";
import {uuid} from './uuid.ts'

type SwapAction = {
    method: () => "GET" | "POST"
    path: () => string
    body: () => unknown
    position: () => number
    withMethod: (method: "GET" | "POST") => SwapAction
    withPath: (path: string) => SwapAction
    withBody: (body: unknown) => SwapAction
    play: (update: boolean) => Promise<void>
}

let nextPosition = 0
const record = {} as Record<string, SwapAction>

function find(id: string): false | SwapAction {
    return record[id] ?? false
}

function swap(view: View<unknown>): SwapAction {
    let swapMethod = 'GET' as "GET" | "POST"
    let swapPath = location.pathname
    let swapBody: unknown
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
        withMethod(method: "GET" | "POST") {
            swapMethod = method
            return this
        },
        withPath(path: string) {
            swapPath = path
            return this
        },
        withBody(body: unknown) {
            swapBody = body
            return this
        },

        async play(update: boolean) {
            const response = await fetch(swapPath, {
                method: swapMethod,
                headers: {Accept: "application/json"},
                body: swapBody as BodyInit
            });

            const json = await response.json();

            view.data = json.data
            view.name = json.name;
            view.error = json.error;

            if (update) {
                const id = uuid()
                record[id] = this
                window.history.pushState(id, "", response.url);
            }
        }
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