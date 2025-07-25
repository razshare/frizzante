<script lang="ts" module>
    import type { SvelteComponent } from "svelte"
    let PreviousComponent = $state(false) as false | SvelteComponent
    let PreviousProperties = $state({}) as Record<string, unknown>
</script>

<script lang="ts">
    let { from, properties } = $props()
    PreviousProperties = properties
    from.then(function next(view: SvelteComponent) {
        PreviousComponent = view
    })
</script>

{#await from}
    {#if PreviousComponent}
        <PreviousComponent.default {...PreviousProperties} />
    {/if}
{:then Component}
    <Component.default {...properties} />
{/await}
