<script lang="ts">
    import Layout from "$lib/components/Layout.svelte"
    import { href } from "$lib/scripts/core/href.ts"
    import { fade } from "svelte/transition"
</script>

<Layout title="Welcome">
    <main
        in:fade={{ duration: 300 }}
        class="min-h-screen flex items-center justify-center p-4"
    >
        <div class="max-w-4xl w-full space-y-8">
            <div class="text-center space-y-4">
                <h1
                    class="text-6xl md:text-8xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent inline-block leading-tight mb-0"
                >
                    Frizzante
                </h1>
                <p class="text-xl text-base-content/60">
                    Modern Go + Svelte Framework
                </p>
            </div>

            <div class="grid md:grid-cols-3 gap-4">
                <div
                    class="card bg-base-200/50 backdrop-blur border border-base-300/70"
                >
                    <div class="card-body">
                        <h2 class="card-title text-lg text-secondary">
                            Lightning Fast
                        </h2>
                        <p class="text-base-content/70 text-base">
                            Built with Go for exceptional performance
                        </p>
                    </div>
                </div>

                <div
                    class="card bg-base-200/50 backdrop-blur border border-base-300/70"
                >
                    <div class="card-body">
                        <h2 class="card-title text-lg text-secondary">
                            Reactive UI
                        </h2>
                        <p class="text-base-content/70 text-base">
                            Powered by Svelte for smooth interfaces
                        </p>
                    </div>
                </div>

                <div
                    class="card bg-base-200/50 backdrop-blur border border-base-300/70"
                >
                    <div class="card-body">
                        <h2 class="card-title text-lg text-secondary">
                            Developer First
                        </h2>
                        <p class="text-base-content/70 text-base">
                            Opinionated choices that streamline workflow
                        </p>
                    </div>
                </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-4 justify-center">
                <a class="btn btn-primary btn-lg" {...href("/todos")}>
                    Show Todos
                    <span class="text-lg">→</span>
                </a>
                <a
                    class="btn btn-outline btn-lg"
                    href="https://razshare.github.io/frizzante-docs/guides/get-started"
                    target="_blank"
                >
                    Documentation
                </a>
            </div>
        </div>
    </main>
</Layout>
