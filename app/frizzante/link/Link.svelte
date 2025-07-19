<script lang="ts">
    import type { Snippet } from "svelte"
    import { href } from "$frizzante/core/scripts/href.ts"

    type Props = {
        href: string
        children: Snippet
        error?: Snippet<[Error]>
        pending?: Snippet
        class?: string
        style?: string
    }
    let {
        href: path,
        children,
        error: onError,
        pending: onPending,
        class: cls,
        style,
    }: Props = $props()

    let pending: boolean = $state(false)
    let error: false | Error = $state(false)

    let options = $derived.by(function run() {
        const out = href(path)

        return {
            href: out.href,
            onclick(event: MouseEvent) {
                pending = true
                out.onclick(event)
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

<a {...options} class={cls} {style}>
    {#if error && onError}
        {@render onError(error)}
    {:else if pending && onPending}
        {@render onPending()}
    {:else}
        {@render children()}
    {/if}
</a>
