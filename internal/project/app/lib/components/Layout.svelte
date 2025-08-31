<style>
    .content {
        position: fixed;
        left: 0;
        right: 0;
        top: 0;
        bottom: 0;
        background: #1e1e2e;
        color: #ed6f49;
        display: grid;
        justify-content: center;
        align-content: center;
        font-family: "Noto Sans Gothic", serif;
        text-align: center;
    }
</style>

<script lang="ts">
    import { getContext, type Snippet } from "svelte"
    import type { View } from "$lib/scripts/core/types.ts"

    type Props = {
        title?: string
        children: Snippet
    }

    const view = getContext("view") as View<unknown>

    let { title = view.name, children }: Props = $props()
</script>

<svelte:head>
    <meta charset="UTF-8" />
    <meta
        name="viewport"
        content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"
    />
    <title>{title}</title>
</svelte:head>

<div class="content">
    {@render children()}
</div>
