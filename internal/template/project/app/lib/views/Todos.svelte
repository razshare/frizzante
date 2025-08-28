<style>
    ol {
        padding: 1rem;
        border-radius: 0.3rem;
        background: rgba(0, 0, 0, 0.3);
        list-style-type: none;
        min-width: 400px;
        text-align: start;
    }

    input {
        background: transparent;
        border: 0;
        border-bottom: 1px solid cadetblue;
        padding: 0.3rem;
        border-radius: 0;
        color: cadetblue;
    }

    input:focus {
        background-color: rgba(0, 0, 0, 0.1);
        outline: none;
    }
</style>

<script lang="ts">
    import Layout from "$lib/components/Layout.svelte"
    import { href } from "$frizzante/core/scripts/href.ts"
    import Form from "$frizzante/forms/components/Form.svelte"

    type Todo = {
        Checked: boolean
        Description: string
    }

    type Props = {
        todos: Todo[]
        error: string
    }

    let { todos, error }: Props = $props()
</script>

<Layout title="Todos">
    <ol>
        {#each todos as todo, index (index)}
            <li>
                <Form action="/remove">
                    <input type="hidden" name="index" value={index} />
                    <button class="link">[Remove]</button>
                </Form>
                {#if todo.Checked}
                    <Form action="/uncheck">
                        <input type="hidden" name="index" value={index} />
                        <button class="link">
                            <!---->
                            (x) {todo.Description}
                            <!---->
                        </button>
                    </Form>
                {:else}
                    <Form action="/check">
                        <input type="hidden" name="index" value={index} />
                        <button class="link">
                            <!---->
                            (&nbsp;&nbsp;) {todo.Description}
                            <!---->
                        </button>
                    </Form>
                {/if}
            </li>
        {/each}
    </ol>
    <Form action="/add">
        <span class="link">Description</span>
        <input type="text" value="" name="description" />
        <button class="link" type="submit">Add +</button>
    </Form>

    {#if error}
        <br />
        <span class="error">{error}</span>
    {/if}

    <br />
    <a class="link" {...href("/")}>&lt; Back</a>
</Layout>
