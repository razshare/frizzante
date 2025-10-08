<script lang="ts">
    import Icon from "$lib/components/icons/Icon.svelte"
    import Layout from "$lib/components/Layout.svelte"
    import { action } from "$lib/scripts/core/action.ts"
    import { href } from "$lib/scripts/core/href.ts"
    import { mdiArrowLeft, mdiCheckCircleOutline, mdiCircleOutline, mdiClose, mdiPlus } from "@mdi/js"
    import { slide } from "svelte/transition"
    import type { Props, sessions } from "$gen/types/main/lib/routes/todos/Props"

    let { items = [], error }: Props = $props()

    let unchecked = $derived.by(function count(): number {
        let value = 0
        for (const todo of items) {
            if (!todo.checked) {
                value++
            }
        }
        return value
    })
</script>

<Layout title="Todos">
    <div class="w-full min-w-[450px] max-w-2xl">
        <div class="text-center">
            {@render Description()}
        </div>
        <div class="card-body relative p-6">
            {@render AddTodoForm()}
            <div class="divider"></div>
            {@render ShowTodosList(items)}
            {@render BackButton()}
        </div>
    </div>
</Layout>

{#snippet Description()}
    <h1 class="text-3xl">Daily Tasks</h1>
    <p class="text-lg text-base-content/60">Organize and track your daily activities</p>
{/snippet}

{#snippet AddTodoForm()}
    <form {...action("/add")} class="flex">
        <input
            type="text"
            name="description"
            placeholder="Add a new task..."
            class="input bg-base-100/ text-lg w-full"
        />
        <div class="pt-4"></div>
        <button type="submit" class="btn btn-ghost text-lg">
            <Icon path={mdiPlus} size="20" />
            <span>Add</span>
        </button>
    </form>

    {#if error}
        <div class="pt-4"></div>
        <div in:slide out:slide class="alert alert-error">
            <span>{error}</span>
        </div>
    {/if}
{/snippet}

{#snippet ShowTodosList(items: sessions.Todo[])}
    {#if items.length > 0}
        {#each items as todo, index (index)}
            <div in:slide out:slide class="flex w-full text-base-content/80">
                {@render ToggleTodoButton(todo, index)}
                {@render RemoveTodoButton(index)}
            </div>
        {/each}
        {@render CountUncheckedTodos()}
    {:else}
        {@render NoTodosFound()}
    {/if}
{/snippet}

{#snippet NoTodosFound()}
    <div in:slide out:slide class="text-center text-base-content/50 text-lg">
        <span>No tasks yet. Add one above to get started!</span>
    </div>
{/snippet}

{#snippet ToggleTodoButton(todo: sessions.Todo, index: number)}
    {@const aria = todo.checked ? "Uncheck" : "Check"}
    {@const value = todo.checked ? "0" : "1"}
    {@const icon = todo.checked ? mdiCheckCircleOutline : mdiCircleOutline}
    <form {...action("/toggle")} class="grow content-center">
        <input type="hidden" name="index" value={index} />
        <input type="hidden" name="value" {value} />
        <button
            type="submit"
            class="w-full flex cursor-pointer"
            class:line-through={todo.checked}
            class:text-base-content={todo.checked}
            class:opacity-50={todo.checked}
            aria-label={aria}
        >
            <Icon path={icon} />
            <div class="pr-4"></div>
            <span>{todo.description}</span>
        </button>
    </form>
{/snippet}

{#snippet RemoveTodoButton(index: number)}
    <form {...action("/remove")}>
        <input type="hidden" name="index" value={index} />
        <button
            type="submit"
            class="btn btn-ghost btn-sm btn-square hover:text-error hover:bg-error/20 transition-colors"
            aria-label="Delete"
        >
            <Icon path={mdiClose} size="18" />
        </button>
    </form>
{/snippet}

{#snippet CountUncheckedTodos()}
    <div in:slide out:slide class="text-lg text-base-content/50 text-center">
        <span>{unchecked} tasks remaining</span>
    </div>
{/snippet}

{#snippet BackButton()}
    <div class="pt-4"></div>
    <a class="btn btn-neutral text-lg" {...href("/")}>
        <Icon path={mdiArrowLeft} size="18" />
        <span>Back</span>
    </a>
{/snippet}
