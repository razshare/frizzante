<script lang="ts">
    import Layout from '$lib/components/layout.svelte'
    import {getContext} from "svelte";
    import type {ServerContext} from "$frizzante/types.ts";
    import Router from "$frizzante/components/Router.svelte";

    const server = getContext("server") as ServerContext<any>
</script>

<svelte:head>
    <title>Welcome</title>
</svelte:head>

<Router/>
<Layout>
    <h1>Hello {server.data.name}.</h1>
</Layout>