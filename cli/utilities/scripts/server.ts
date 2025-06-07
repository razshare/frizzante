import {render as _render} from "svelte/server";
import ServerView from "$lib/utilities/components/ServerView.svelte";
export async function render(props: never) {
    return _render(ServerView, {props});
}