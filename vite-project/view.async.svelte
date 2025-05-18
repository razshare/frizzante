<script module>
    let PreviousComponent = $state(false)
</script>
<script>
    let {from, server = $bindable()} = $props()
    from.then((view) => {
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