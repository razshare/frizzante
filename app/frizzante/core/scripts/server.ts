import { render as _render } from "svelte/server"
import ServerView from "$frizzante/core/components/ServerView.svelte"
export async function render(props: unknown) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    return _render(ServerView, { props })
}
