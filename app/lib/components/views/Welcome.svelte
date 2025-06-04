<script lang="ts">
    import Layout from '../Layout.svelte'
    import {getContext} from "svelte";
    import type {View} from "frz/.generated/utilities/types.ts";

    const view = getContext("view") as View<{ name: string }>
</script>

<svelte:head>
    <title>Welcome</title>
</svelte:head>

<Layout>
    <h1>Hello {view.data.name}.</h1>
</Layout>