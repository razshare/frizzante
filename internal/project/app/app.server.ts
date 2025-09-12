import { render as _render } from "svelte/server"
import App from "./app.server.svelte"
export async function render(props: unknown) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    return _render(App, { props })
}
