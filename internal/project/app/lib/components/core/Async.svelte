<script lang="ts" module>
    import type { SvelteComponent } from "svelte"
    let PreviousComponent = $state(false) as false | SvelteComponent
    let PreviousProperties = $state({}) as Record<string, unknown>
</script>

<script lang="ts">
    type Props = {
        from: Promise<SvelteComponent>
        properties: Record<string, unknown>
    }
    let { from, properties }: Props = $props()
    PreviousProperties = properties
    from.then(function next(component: SvelteComponent) {
        PreviousComponent = component
    })
</script>

{#await from}
    {#if PreviousComponent}
        <PreviousComponent.default {...PreviousProperties} />
    {/if}
{:then Component}
    <Component.default {...properties} />
{/await}
