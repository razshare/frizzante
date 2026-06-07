<script lang="ts">
    import Icon from "$lib/components/icons/icon.svelte"
    import Layout from "$lib/components/layout.svelte"
    import Logo from "$lib/components/logo.svelte"
    import { href } from "$lib/scripts/core/href.ts"
    import { mdiArrowRight } from "@mdi/js"
</script>

<Layout title="Welcome">
    <Logo />
    {@render Description()}
    <div class="pt-6"></div>
    <div class="flex justify-center gap-2 relative">
        {@render TodosButton()}
        {@render DocumentationButton()}
    </div>
</Layout>

{#snippet Description()}
    <p>Modern Go + Svelte Framework</p>
{/snippet}

{#snippet TodosButton()}
    <a {...href("/todos")}>
        <button>
            <span>Show Todos</span>
            <Icon path={mdiArrowRight} />
        </button>
    </a>
{/snippet}

{#snippet DocumentationButton()}
    <a href="https://razshare.github.io/frizzante-docs/guides/get-started" target="_blank">
        <button>Documentation</button>
    </a>
{/snippet}
