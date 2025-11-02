<script lang="ts">
    import type { Snippet } from "svelte"
    import { href } from "$lib/scripts/core/href.ts"

    type Props = {
        href: string
        children: Snippet<[{ pending: boolean; error: false | Error }]>
        class?: string
        style?: string
    }
    let { href: path, children, class: cls, style }: Props = $props()

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
    {@render children({ pending, error })}
</a>
