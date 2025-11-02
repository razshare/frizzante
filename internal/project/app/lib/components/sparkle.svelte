<script lang="ts" module>
    const frames: string[] = ["⯌", "⯏", "⯎", "⯍"]
</script>

<script lang="ts">
    import { onMount, type Snippet } from "svelte"

    type Props = {
        children: Snippet<[string]>
        /**
         * Refresh rate in milliseconds
         */
        refresh?: number
    }

    let { children, refresh: update = 500 }: Props = $props()

    let index: number = $state(0)
    let value: string = $derived(frames[index])
    onMount(function start() {
        const timeout = setInterval(function run() {
            index = (index + 1) % frames.length
        }, update)

        return function stop() {
            clearInterval(timeout)
        }
    })
</script>

{@render children(value)}
