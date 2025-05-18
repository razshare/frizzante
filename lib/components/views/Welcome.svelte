<script lang="ts">
    import Layout from '$lib/components/Layout.svelte'

    type Props = {
        server: ServerProperties<{
            name: string
        }>
    }

    let {server = $bindable()}:Props = $props()
</script>

<svelte:head>
    <title>Welcome</title>
</svelte:head>

<Layout>
    <h1>Hello {server.data.name}.</h1>
</Layout>