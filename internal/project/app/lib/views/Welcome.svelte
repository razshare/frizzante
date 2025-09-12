<script lang="ts">
    import Icon from "$lib/components/icons/Icon.svelte"
    import Layout from "$lib/components/Layout.svelte"
    import { href } from "$lib/scripts/core/href.ts"
    import { mdiChevronRight } from "@mdi/js"
    import { slide } from "svelte/transition"
</script>

<Layout title="Welcome">
    <div in:slide class="card bg-neutral text-neutral-content w-96">
        <div class="card-body items-center text-center">
            <h2 class="card-title">Welcome to Frizzante!</h2>
            <div>
                Frizzante is an opinionated web server framework written in Go
                that uses Svelte to render web pages. Read more details <a
                    class="link text-info underline"
                    href="https://razshare.github.io/frizzante-docs/guides/get-started"
                    >here</a
                >.
            </div>
            <div class="card-actions justify-end">
                <a class="btn" {...href("/todos")}>
                    <span>Show todos</span>
                    <Icon path={mdiChevronRight} />
                </a>
            </div>
        </div>
    </div>
</Layout>
