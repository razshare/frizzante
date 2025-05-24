import type {ServerContext} from "$frizzante/types.ts";

type Modifier = "push" | "back" | "forward"

type SwapPayload = {
    modifier: Modifier
    method: string
    path: string
    body?: any
    counter?: number
}

let counter = 0
let started = false

export async function swap(server: ServerContext<any>, payload: SwapPayload): Promise<void> {
    const {method, body, path, modifier} = payload
    const response = await fetch(path, {
        method,
        headers: {Accept: "application/json"},
        body
    });
    const json = await response.json();
    server.data = json.data
    server.view = json.view;

    const search = response.url.split('?', 2)[1] ?? ''
    if ("push" === modifier) {
        counter++
        if ('' !== search) {
            window.history.pushState({method, body, path, modifier, counter}, "", `${path}?${search}`);
            return
        }
        window.history.pushState({method, body, path, modifier, counter}, "", path);
    }
}

export function route(server: ServerContext<any>): void {
    if (started) {
        return
    }

    const listener = async function pop(e: PopStateEvent) {
        e.preventDefault();
        let {method, body, path, modifier, counter: counterLocal} = (e.state ?? {
            method: "GET",
            path: "/",
            modifier: "push",
            counter: 0
        }) as SwapPayload

        counterLocal = counterLocal ?? 0

        if (counterLocal <= counter) {
            counter = counterLocal;
            await swap(server,{modifier: "back", method, path, body})
        } else if (counterLocal > counter) {
            counter = counterLocal;
            await swap(server,{modifier: "forward", method, path, body})
        } else {
            await swap(server,{modifier: "push", method, path, body})
        }
    }
    window.addEventListener("popstate", listener);
    started = true
}
