<script lang="ts" module>
    let PreviousComponent = $state(false)
</script>
<script lang="ts">
    let {from} = $props()
    from.then((view:never) => {
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