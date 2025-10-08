<script lang="ts">
    import type { Snippet } from "svelte"
    import { getContext } from "svelte"
    import type { FormContext } from "$lib/scripts/core/form.svelte.ts"

    type Props = {
        children?: Snippet<[{ errors: string[] }]>
        class?: string
        errorClass?: string
    }

    let { children, class: cls, errorClass }: Props = $props()

    // Get field context
    const fieldContext = getContext("field") as {
        name: string
        form: FormContext
    } | undefined

    if (!fieldContext) {
        throw new Error("FieldErrors component must be used within a Field component")
    }

    const { name, form } = fieldContext

    const errors = $derived(form.getFieldErrors(name as never))
    const hasErrors = $derived(errors.length > 0)
    const errorId = `field-${name}-error`
</script>

{#if hasErrors}
    <div id={errorId} class={cls} role="alert" aria-live="polite">
        {#if children}
            {@render children({ errors })}
        {:else}
            {#each errors as error}
                <div class={errorClass || "text-error text-sm mt-1"}>{error}</div>
            {/each}
        {/if}
    </div>
{/if}
