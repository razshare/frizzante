<script lang="ts">
    import Layout from '$lib/components/Layout.svelte'
    import {getContext} from "svelte";
    import type {ServerContext} from "$frizzante/types.ts";
    let server = getContext("server") as ServerContext<any>
</script>

<svelte:head>
    <title>Welcome</title>
</svelte:head>

<Layout>
    <h1>Hello {server.data.name}.</h1>
</Layout>