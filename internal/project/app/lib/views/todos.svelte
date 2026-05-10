<style lang="scss">
  .todo-list {
    padding: 1rem;
    max-height: 300px;
    overflow-y: auto;
    .item {
      width: 100%;
      display: grid;
      grid-template-columns: 1fr auto;
      form {
        text-align: start;
        label {
          button {
            display: none;
          }
        }
      }
    }
  }
</style>

<script lang="ts">
    import type { Props, sessions } from "$gen/types/main/lib/routes/todos/props"
    import Icon from "$lib/components/icons/icon.svelte"
    import Layout from "$lib/components/layout.svelte"
    import { action } from "$lib/scripts/core/action.ts"
    import { href } from "$lib/scripts/core/href.ts"
    import { mdiArrowLeft, mdiCheckCircleOutline, mdiCircleOutline, mdiClose, mdiPlus } from "@mdi/js"
    import { slide } from "svelte/transition"
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
    {@render Description()}
    {@render AddTodoForm()}
    {@render TodoList(items)}
    {@render BackButton()}
</Layout>

{#snippet Description()}
    <h1>Daily Tasks</h1>
    <p>Organize and track your daily activities.</p>
{/snippet}

{#snippet AddTodoForm()}
    <form method="POST" {...action("/add")}>
        <input type="text" name="description" placeholder="Add a new task..." />
        <br/>
        <button type="submit">
            <Icon path={mdiPlus} />
            <span>Add</span>
        </button>
    </form>
    {#if error}
        <br />
        <div transition:slide>
            <span>{error}</span>
        </div>
    {/if}
{/snippet}

{#snippet TodoList(items: sessions.Todo[])}
    {#if items.length > 0}
        <div class="todo-list">
            {#each items as todo, index (index)}
                <div transition:slide class="item">
                    {@render ToggleTodoButton(todo, index)}
                    {@render RemoveTodoButton(index)}
                </div>
            {/each}
        </div>
        {@render CountUncheckedTodos()}
    {:else}
        <div class="todo-list">
            {@render NoTodosFound()}
        </div>
    {/if}
{/snippet}

{#snippet NoTodosFound()}
    <span>No tasks yet. Add one above to get started!</span>
{/snippet}

{#snippet ToggleTodoButton(todo: sessions.Todo, index: number)}
    {@const aria = todo.checked ? "Uncheck" : "Check"}
    {@const value = todo.checked ? "0" : "1"}
    {@const icon = todo.checked ? mdiCheckCircleOutline : mdiCircleOutline}
    <form method="POST" {...action("/toggle")}>
        <input type="hidden" name="index" value={index} />
        <input type="hidden" name="value" {value} />
        <label
        // type="submit"
        // class="w-full flex cursor-pointer"
        // class:line-through={todo.checked}
        // class:text-base-content={todo.checked}
        // class:opacity-50={todo.checked}
        aria-label={aria}
        >
        <Icon path={icon} />
        {#if todo.checked}
            <strike>
                <span>{todo.description}</span>
            </strike>
        {:else}
            <span>{todo.description}</span>
        {/if}
        <button aria-label="toggle"></button>
        </label>
    </form>
{/snippet}

{#snippet RemoveTodoButton(index: number)}
    <form method="POST" {...action("/remove")}>
        <input type="hidden" name="index" value={index} />
        <label aria-label="Delete">
            <Icon path={mdiClose} />
            <button aria-label="remove"></button>
        </label>
    </form>
{/snippet}

{#snippet CountUncheckedTodos()}
    <span>{unchecked} tasks remaining</span>
{/snippet}

{#snippet BackButton()}
    <br />
    <a {...href("/")}>
        <Icon path={mdiArrowLeft} />
        <span>Back</span>
    </a>
{/snippet}
