<script lang="ts">
    import Icon from "$lib/components/icons/Icon.svelte"
    import Layout from "$lib/components/Layout.svelte"
    import { action } from "$lib/scripts/core/action.ts"
    import { href } from "$lib/scripts/core/href.ts"
    import {
        mdiCheck,
        mdiChevronLeft,
        mdiCloseCircle,
        mdiInformationVariantCircle,
        mdiMinusCircle,
        mdiPlusCircle,
    } from "@mdi/js"
    import { slide } from "svelte/transition"

    type Todo = {
        Checked: boolean
        Description: string
    }

    type Props = {
        todos: Todo[]
        mode: number
        error: string
    }

    let { todos, mode, error }: Props = $props()
</script>

{#snippet Toggle()}
    <a class="btn btn-ghost absolute left-1 top-1" {...href("/")}>
        <Icon path={mdiChevronLeft} />
        <span>Welcome</span>
    </a>
    <div class="absolute right-1 top-1">
        <a class="btn btn-ghost" {...href("/mode?value=remove")}>
            <span>Remove</span>
            <Icon path={mdiMinusCircle} />
        </a>
        <a class="btn btn-ghost" {...href("/mode?value=add")}>
            <span>Add</span>
            <Icon path={mdiPlusCircle} />
        </a>
    </div>
    <div class="pt-8"></div>
    <div class="card-body items-center text-start">
        <h2 class="card-title text-3xl">My Todos</h2>
        <div class="badge badge-info p-4">
            <Icon path={mdiInformationVariantCircle} />
            <span>Select an item in order to check or uncheck it.</span>
        </div>
        <div class="card-actions">
            {#each todos as todo, index (index)}
                {@const url = todo.Checked ? "/uncheck" : "/check"}
                <form in:slide class="w-full" {...action(url)}>
                    <input type="hidden" name="index" value={index} />
                    <button class="relative btn btn-ghost w-full" type="submit">
                        {#if todo.Checked}
                            <div
                                in:slide={{ axis: "x" }}
                                out:slide={{ axis: "x" }}
                                class="absolute left-1"
                            >
                                <Icon path={mdiCheck} />
                            </div>
                        {/if}
                        <div class="pr-2"></div>
                        <span>{todo.Description}</span>
                    </button>
                </form>
            {/each}
        </div>
    </div>
{/snippet}

{#snippet Add()}
    <a
        class="btn btn-ghost absolute left-1 top-1"
        {...href("/mode?value=toggle")}
    >
        <Icon path={mdiChevronLeft} />
        <span>Toggle</span>
    </a>
    <div class="pt-8"></div>
    <div class="card-body items-center text-start">
        <h2 class="card-title text-3xl">Add a new Todo</h2>
        <div class="card-actions">
            <form {...action("/add")}>
                <div class="join">
                    <input
                        class="input join-item"
                        name="description"
                        placeholder="Description"
                    />
                    <button class="btn join-item rounded-r-full" type="submit">
                        <span>Add</span>
                        <Icon path={mdiPlusCircle} />
                    </button>
                </div>
            </form>
        </div>
    </div>
{/snippet}

{#snippet Remove()}
    <a
        class="btn btn-ghost absolute left-1 top-1"
        {...href("/mode?value=toggle")}
    >
        <Icon path={mdiChevronLeft} />
        <span>Toggle</span>
    </a>
    <div class="pt-8"></div>
    <div class="card-body items-center text-start">
        <h2 class="card-title text-3xl">Remove a Todo</h2>
        <div class="card-actions">
            {#each todos as todo, index (index)}
                {@const url = "/remove"}
                <form in:slide class="w-full" {...action(url)}>
                    <input type="hidden" name="index" value={index} />
                    <button class="relative btn btn-ghost w-full" type="submit">
                        {#if todo.Checked}
                            <div
                                in:slide={{ axis: "x" }}
                                out:slide={{ axis: "x" }}
                                class="absolute left-1"
                            >
                                <Icon path={mdiCheck} />
                            </div>
                        {/if}
                        <div class="pr-2"></div>
                        <span>{todo.Description}</span>
                    </button>
                </form>
            {/each}
        </div>
    </div>
{/snippet}

<Layout title="Todos">
    <div in:slide class="card bg-neutral text-neutral-content w-96 relative">
        {#if mode === 0}
            {@render Toggle()}
        {:else if mode === 1}
            {@render Remove()}
        {:else if mode === 2}
            {@render Add()}
        {/if}
        {#if error}
            <div class="pt-4"></div>
            <div class="flex text-error p-4">
                <Icon path={mdiCloseCircle} />
                <div class="pr-2"></div>
                <span>{error}</span>
            </div>
        {/if}
    </div>
</Layout>
