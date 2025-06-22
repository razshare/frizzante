<style global>
    .content {
        position: fixed;
        left: 0;
        right: 0;
        top: 0;
        bottom: 0;
        background: #1e1e2e;
        color: cadetblue;
        display: grid;
        justify-content: center;
        align-content: center;
        font-family: "Noto Sans Gothic", serif;
    }
</style>

<script lang="ts">
    const { children } = $props()
</script>

<svelte:head>
    <meta charset="UTF-8" />
    <meta
        name="viewport"
        content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"
    />
</svelte:head>

<div class="content">
    {@render children()}
</div>
