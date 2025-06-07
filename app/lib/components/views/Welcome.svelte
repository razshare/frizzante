<script lang="ts">
    import {getContext} from "svelte";
    import type {View} from "$frz/types";
    import Layout from "$lib/components/Layout.svelte";

    const view = getContext("view") as View<{ name: string }>
</script>

<svelte:head>
    <title>Welcome</title>
</svelte:head>

<Layout>
    <h1>Hello {view.data.name}.</h1>
</Layout>