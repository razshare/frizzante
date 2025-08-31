<script lang="ts">
    import Layout from "$lib/components/Layout.svelte"
    import {href} from "$lib/scripts/core/href.ts";
</script>

<Layout title="Welcome">
    <h1>Welcome to Frizzante.</h1>
    <a {...href("/todos")}>Show todos</a>
</Layout>
