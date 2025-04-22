<script module>
    let PreviousComponent = $state(false)
</script>
<script>
    let {from} = $props()
    from.then((view) => {
        PreviousComponent = view
    })
</script>

{#await from}
    {#if PreviousComponent}
        <PreviousComponent.default/>
    {/if}
{:then Component}
    <Component.default/>
{/await}