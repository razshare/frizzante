<script lang="ts" module>
    import type { SvelteComponent } from "svelte"
    let PreviousComponent = $state(false) as false | SvelteComponent
</script>

<script lang="ts">
    let { from } = $props()
    from.then(function next(view: SvelteComponent) {
        PreviousComponent = view
    })
</script>

{#await from}
    {#if PreviousComponent}
        <PreviousComponent.default />
    {/if}
{:then Component}
    <Component.default />
{/await}
