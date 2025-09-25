import type {SvelteComponent} from "svelte";
import Welcome from "$lib/views/Welcome.svelte";
import Todos from "$lib/views/Todos.svelte";

export const views:Record<string, SvelteComponent> = {
    Welcome: Welcome as SvelteComponent,
    Todos: Todos as SvelteComponent,
}
