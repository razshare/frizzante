<script lang="ts">
    import Layout from "$lib/components/Layout.svelte"
    import { href } from "$lib/scripts/core/href.ts"
    import Logo from "$lib/components/Logo.svelte"
    import { mdiArrowRight } from "@mdi/js"
    import Icon from "$lib/components/icons/Icon.svelte"
    import Sparkle from "$lib/components/Sparkle.svelte"
    import { scale } from "svelte/transition"
</script>

<Layout title="Welcome">
    <Logo />
    {@render Description()}
    <div class="pt-6"></div>
    <div class="flex justify-center gap-2 relative">
        {@render BackgroundEffect()}
        {@render TodosButton()}
        {@render DocumentationButton()}
    </div>
</Layout>

{#snippet Description()}
    <div class="pt-6"></div>
    <p class="text-xl text-base-content/60">Modern Go + Svelte Framework</p>
{/snippet}

{#snippet BackgroundEffect()}
    {@const color = "text-secondary"}
    <Sparkle refresh={500}>
        {#snippet children(frame: string)}
            <span in:scale={{ duration: 500 * 10 }} class="{color} absolute left-32 -top-10">{frame}</span>
        {/snippet}
    </Sparkle>
    <Sparkle refresh={300}>
        {#snippet children(frame: string)}
            <span in:scale={{ duration: 300 * 10 }} class="{color} absolute left-24 -top-52">{frame}</span>
        {/snippet}
    </Sparkle>
    <Sparkle refresh={1200}>
        {#snippet children(frame: string)}
            <span in:scale={{ duration: 1200 * 10 }} class="{color} absolute right-4 -top-40">{frame}</span>
        {/snippet}
    </Sparkle>
    <Sparkle refresh={570}>
        {#snippet children(frame: string)}
            <span in:scale={{ duration: 570 * 10 }} class="{color} absolute right-20 -top-14">{frame}</span>
        {/snippet}
    </Sparkle>
    <Sparkle refresh={330}>
        {#snippet children(frame: string)}
            <span in:scale={{ duration: 330 * 10 }} class="{color} absolute right-40 -top-56">{frame}</span>
        {/snippet}
    </Sparkle>
    <Sparkle refresh={400}>
        {#snippet children(frame: string)}
            <span in:scale={{ duration: 400 * 10 }} class="{color} absolute left-40 -top-56">{frame}</span>
        {/snippet}
    </Sparkle>
    <Sparkle refresh={490}>
        {#snippet children(frame: string)}
            <span in:scale={{ duration: 490 * 10 }} class="{color} absolute -right-10 -top-16">{frame}</span>
        {/snippet}
    </Sparkle>
    <Sparkle refresh={350}>
        {#snippet children(frame: string)}
            <span in:scale={{ duration: 350 * 10 }} class="{color} absolute -left-10 -top-16">{frame}</span>
        {/snippet}
    </Sparkle>
    <Sparkle refresh={560}>
        {#snippet children(frame: string)}
            <span in:scale={{ duration: 560 * 10 }} class="{color} absolute text-center -top-60">{frame}</span>
        {/snippet}
    </Sparkle>
{/snippet}

{#snippet TodosButton()}
    <a class="btn btn-primary btn-lg" {...href("/todos")}>
        <span>Show Todos</span>
        <Icon path={mdiArrowRight} size="18" />
    </a>
{/snippet}

{#snippet DocumentationButton()}
    <a
        class="btn btn-secondary btn-lg"
        href="https://razshare.github.io/frizzante-docs/guides/get-started"
        target="_blank"
    >
        <span>Documentation</span>
    </a>
{/snippet}
