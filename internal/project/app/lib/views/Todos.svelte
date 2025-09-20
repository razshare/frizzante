<script lang="ts">
    import Icon from "$lib/components/icons/Icon.svelte"
    import Layout from "$lib/components/Layout.svelte"
    import { action } from "$lib/scripts/core/action.ts"
    import { href } from "$lib/scripts/core/href.ts"
    import { mdiArrowLeft, mdiCheck, mdiDelete, mdiPlus } from "@mdi/js"
    import { fade, slide } from "svelte/transition"
    import type {Props} from "$lib/types/gen/main/lib/routes/handlers/todos/Props";

    let { todos = [], error }: Props = $props()
</script>

<Layout title="Todos">
    <main
        in:fade={{ duration: 300 }}
        class="min-h-screen flex items-center justify-center p-4"
    >
        <div class="w-full min-w-[450px] max-w-2xl space-y-8">
            <div class="text-center space-y-2 mb-5">
                <a
                    class="btn btn-ghost text-lg gap-1 absolute left-4 top-4"
                    {...href("/")}
                >
                    <Icon path={mdiArrowLeft} size="18" />
                    Back
                </a>
                <h1
                    class="text-4xl md:text-5xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent inline-block leading-tight mb-0"
                >
                    My Tasks
                </h1>
                <p class="text-lg text-base-content/60">
                    Organize and track your daily activities
                </p>
            </div>

            <div
                class="card bg-base-200/50 backdrop-blur border border-base-300/70"
            >
                <div class="card-body p-6 space-y-4">
                    <form {...action("/add")} class="flex gap-2 mb-0">
                        <input
                            type="text"
                            name="description"
                            placeholder="Add a new task..."
                            class="input input-bordered flex-1 bg-base-100/ text-lg"
                            required
                        />
                        <button type="submit" class="btn btn-primary text-lg">
                            <Icon path={mdiPlus} size="20" />
                            Add
                        </button>
                    </form>

                    {#if error}
                        <div in:slide class="alert alert-error">
                            <span>{error}</span>
                        </div>
                    {/if}

                    <div class="divider my-0.5"></div>

                    {#if todos.length === 0}
                        <div
                            class="text-center py-4 text-base-content/50 text-lg"
                        >
                            No tasks yet. Add one above to get started!
                        </div>
                    {:else}
                        <div class="space-y-2">
                            {#each todos as todo, index (index)}
                                <div
                                    in:slide
                                    class="group flex items-center gap-3 p-3 rounded-lg bg-base-100 border border-base-300/70 hover:border-primary/50 hover:shadow-lg hover:shadow-primary/10 transition-all"
                                >
                                    <form
                                        {...action(
                                            todo.checked
                                                ? "/uncheck"
                                                : "/check",
                                        )}
                                        class="flex-shrink-0"
                                    >
                                        <input
                                            type="hidden"
                                            name="index"
                                            value={index}
                                        />
                                        <button
                                            type="submit"
                                            class="btn btn-ghost btn-sm btn-square"
                                            aria-label={todo.checked
                                                ? "Uncheck"
                                                : "Check"}
                                        >
                                            <div
                                                class={`w-5 h-5 rounded border-2 flex items-center justify-center transition-all ${
                                                    todo.checked
                                                        ? "bg-primary border-primary"
                                                        : "border-base-content/30 hover:border-primary"
                                                }`}
                                            >
                                                {#if todo.checked}
                                                    <Icon
                                                        path={mdiCheck}
                                                        size="14"
                                                    />
                                                {/if}
                                            </div>
                                        </button>
                                    </form>

                                    <span
                                        class={`flex-1 text-lg ${todo.checked ? "line-through text-base-content/50" : ""}`}
                                    >
                                        {todo.description}
                                    </span>

                                    <form
                                        {...action("/remove")}
                                        class="flex-shrink-0"
                                    >
                                        <input
                                            type="hidden"
                                            name="index"
                                            value={index}
                                        />
                                        <button
                                            type="submit"
                                            class="btn btn-ghost btn-sm btn-square text-error hover:bg-error/20 transition-colors"
                                            aria-label="Delete"
                                        >
                                            <Icon path={mdiDelete} size="18" />
                                        </button>
                                    </form>
                                </div>
                            {/each}
                        </div>
                    {/if}

                    {#if todos.length > 0}
                        <div class="text-lg text-base-content/50 text-center">
                            {todos.filter(t => !t.checked).length} of {todos.length}
                            tasks remaining
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </main>
</Layout>
