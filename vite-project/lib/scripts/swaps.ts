import type {ServerContext} from "../types.ts";
import {uuid} from './uuid.ts'

type SwapAction = {
    method: () => "GET" | "POST"
    path: () => string
    body: () => any
    position: () => number
    withMethod: (method: "GET" | "POST") => SwapAction
    withPath: (path: string) => SwapAction
    withBody: (body: any) => SwapAction
    withMerge: (merge: boolean) => SwapAction
    play: (update: boolean) => Promise<void>
}

let nextPosition = 0
let record = {} as Record<string, SwapAction>

function find(id: string): false | SwapAction {
    return record[id] ?? false
}

function swap(server: ServerContext<any>): SwapAction {
    let swapMethod = 'GET' as "GET" | "POST"
    let swapPath = location.pathname
    let swapBody: any
    let swapPosition = nextPosition++
    let swapMerge = false

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
        withBody(body: any) {
            swapBody = body
            return this
        },
        withMerge(merge: boolean) {
            swapMerge = merge
            return this
        },
        async play(update: boolean) {
            const response = await fetch(swapPath, {
                method: swapMethod,
                headers: {Accept: "application/json"},
                body: swapBody
            });

            const json = await response.json();

            if (swapMerge) {
                server.data = {
                    ...server.data,
                    ...json.data,
                }
            } else {
                server.data = json.data
            }

            server.view = json.view;
            server.error = json.error;

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