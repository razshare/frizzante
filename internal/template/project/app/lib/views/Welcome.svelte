<script lang="ts">
    import Layout from "$lib/components/Layout.svelte"
    import Link from "$frizzante/links/components/Link.svelte"
</script>

<Layout title="Welcome">
    <h1>Welcome to Frizzante.</h1>
    <Link class="link" href="/todos">
        <span>Show todos</span>
    </Link>
</Layout>
