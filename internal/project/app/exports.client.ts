import type {SvelteComponent} from "svelte"

export const views:Record<string, Promise<SvelteComponent>> = {
    Welcome: import("$lib/views/Welcome.svelte"),
    Todos: import("$lib/views/Todos.svelte"),
}
