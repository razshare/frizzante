<script lang="ts">
    import Layout from "$lib/components/Layout.svelte"

    type Props = { name: string }
    let { name }: Props = $props()
</script>

<svelte:head>
    <title>Welcome</title>
</svelte:head>

<Layout>
    <h1>Hello {name}.</h1>
</Layout>
