<script lang="ts" module>
    let PreviousComponent = $state(false)
</script>
<script lang="ts">
    let {from, server = $bindable()} = $props()
    from.then((view:any) => {
        PreviousComponent = view
    })
</script>

{#await from}
    {#if PreviousComponent}
        <PreviousComponent.default bind:server/>
    {/if}
{:then Component}
    <Component.default bind:server/>
{/await}