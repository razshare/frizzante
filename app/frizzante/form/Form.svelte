<style>
    form {
        display: inline-block;
    }
</style>

<script lang="ts">
    import { action } from "$frizzante/core/scripts/action.ts"
    import type { Snippet } from "svelte"

    type Props = {
        method?: "GET" | "POST"
        action: string
        children: Snippet
        error?: Snippet<[Error]>
        pending?: Snippet
        class?: string
        style?: string
    }
    let {
        method = "GET",
        action: actionPath,
        children,
        error: onError,
        pending: onPending,
        class: cls,
        style,
    }: Props = $props()
    let pending: boolean = $state(false)
    let error: false | Error = $state(false)
    let options = $derived.by(function run() {
        const out = action(actionPath)

        return {
            action: out.action,
            onsubmit(event: Event) {
                pending = true
                out.onsubmit(event)
                    .then(function run() {
                        pending = false
                    })
                    .catch(function run(errorLocal: Error) {
                        error = errorLocal
                    })
            },
        }
    })
</script>

<form {method} {...options} class={cls} {style}>
    {#if error && onError}
        {@render onError(error)}
    {:else if pending && onPending}
        {@render onPending()}
    {:else}
        {@render children()}
    {/if}
</form>
