import type {ServerContext} from "$frizzante/types.ts";

type Modifier = "push" | "back" | "forward"

let counter = 0;

async function swap(server: ServerContext<any>, id: string, modifier: Modifier, data: any = false): Promise<void> {
    if (!server.ids[id]) {
        return;
    }

    const path = server.ids[id];


    if (false !== data) {
        if ("push" === modifier) {
            window.history.pushState({id, counter: ++counter}, "", path);
        }
        server.id = id;
        return;
    }

    const response = await fetch(path, {
        headers: {Accept: "application/json"},
    });
    const json = await response.json();
    server.data = json.data
    server.ids = json.ids
    server.id = json.id;

    const search = response.url.split('?', 2)[1] ?? ''

    if ("push" === modifier) {
        if ('' !== search) {
            window.history.pushState({id, counter: ++counter}, "", `${path}?${search}`);
            return
        }
        window.history.pushState({id, counter: ++counter}, "", path);
    }
}


export function navigate(server: ServerContext<any>, id: string, data: any = false): Promise<void> {
    return swap(server, id, "push", data);
}

let started = false

export function route(server: ServerContext<any>): void {
    if (started) {
        return
    }

    const listener = async function pop(e: PopStateEvent) {
        e.preventDefault();
        let id = e.state?.id ?? "";

        if ('' === id) {
            for (const idLocal in server.ids) {
                const path = server.ids[idLocal]
                if ('/' === path || '' === path) {
                    id = idLocal
                    break
                }
            }
        }
        const counterLocal = e.state?.counter ?? 0;
        if (counterLocal < counter) {
            counter = counterLocal;
            await swap(server, id, "back");
        } else if (counterLocal > counter) {
            counter = counterLocal;
            await swap(server, id, "forward");
        } else {
            await swap(server, id, "push");
        }
    }
    window.addEventListener("popstate", listener);
    started = true
}
