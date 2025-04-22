<script>
    import Layout from '$lib/components/Layout.svelte'
    import {getContext} from "svelte";

    /**
     * @typedef Data
     * @property {string} name
     */

    /** @type {Data} */
    const data = getContext("data")
</script>

<svelte:head>
    <title>Welcome</title>
</svelte:head>

<Layout>
    <h1>Hello {data.name}.</h1>
</Layout>