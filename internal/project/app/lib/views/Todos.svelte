<script lang="ts">
    import Icon from "$lib/components/icons/Icon.svelte"
    import Layout from "$lib/components/Layout.svelte"
    import { action } from "$lib/scripts/core/action.ts"
    import { href } from "$lib/scripts/core/href.ts"
    import { mdiArrowLeft, mdiCheckCircleOutline, mdiCircleOutline, mdiClose, mdiPlus } from "@mdi/js"
    import { slide } from "svelte/transition"
    import type { Props, Todo} from "$gen/types/main/lib/routes/handlers/todos/Props"

    let { Todos = [], Error }: Props = $props()

    let unchecked = $derived.by(function count(): number {
        let value = 0
        for (const todo of Todos) {
            if (!todo.Checked) {
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
            {@render Add()}
            <div class="divider"></div>
            {#if Todos.length === 0}
                {@render Empty()}
            {:else}
                {#each Todos as todo, index (index)}
                    {@render Todo(todo, index)}
                {/each}
            {/if}
            {#if Todos.length > 0}
                {@render Remaining()}
            {/if}
            {@render Back()}
        </div>
    </div>
</Layout>

{#snippet Description()}
    <h1 class="text-3xl">Daily Tasks</h1>
    <p class="text-lg text-base-content/60">Organize and track your daily activities</p>
{/snippet}

{#snippet Add()}
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

    {#if Error}
        <div class="pt-4"></div>
        <div in:slide out:slide class="alert alert-error">
            <span>{Error}</span>
        </div>
    {/if}
{/snippet}

{#snippet Empty()}
    <div in:slide out:slide class="text-center text-base-content/50 text-lg">
        <span>No tasks yet. Add one above to get started!</span>
    </div>
{/snippet}

{#snippet Todo(todo: Todo, index: number)}
    <div in:slide out:slide class="flex w-full text-base-content/80">
        {@render Toggle(todo, index)}
        {@render Remove(index)}
    </div>
{/snippet}

{#snippet Toggle(todo: Todo, index: number)}
    {@const aria = todo.Checked ? "Uncheck" : "Check"}
    {@const value = todo.Checked ? "0" : "1"}
    {@const icon = todo.Checked ? mdiCheckCircleOutline : mdiCircleOutline}
    <form {...action("/toggle")} class="grow content-center">
        <input type="hidden" name="index" value={index} />
        <input type="hidden" name="value" {value} />
        <button
            type="submit"
            class="w-full flex cursor-pointer"
            class:line-through={todo.Checked}
            class:text-base-content={todo.Checked}
            class:opacity-50={todo.Checked}
            aria-label={aria}
        >
            <Icon path={icon} />
            <div class="pr-4"></div>
            <span>{todo.Description}</span>
        </button>
    </form>
{/snippet}

{#snippet Remove(index: number)}
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

{#snippet Remaining()}
    <div in:slide out:slide class="text-lg text-base-content/50 text-center">
        <span>{unchecked} tasks remaining</span>
    </div>
{/snippet}

{#snippet Back()}
    <div class="pt-4"></div>
    <a class="btn btn-neutral text-lg" {...href("/")}>
        <Icon path={mdiArrowLeft} size="18" />
        <span>Back</span>
    </a>
{/snippet}
