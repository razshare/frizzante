<script lang="ts">
    import type { Snippet } from "svelte"
    import { getContext } from "svelte"
    import type { FormContext } from "$lib/scripts/core/form.svelte.ts"

    type Props = {
        children: Snippet<
            [
                {
                    props: {
                        id: string
                        name: string
                        "aria-invalid": boolean
                        "aria-describedby": string
                        onblur?: (event: FocusEvent) => void
                        oninput?: (event: Event) => void
                    }
                },
            ]
        >
        class?: string
    }

    let { children, class: cls }: Props = $props()

    // Get field context
    const fieldContext = getContext("field") as {
        name: string
        form: FormContext
        validateOnBlur: boolean
        validateOnChange: boolean
    }

    if (!fieldContext) {
        throw new Error("Control component must be used within a Field component")
    }

    const { name, form, validateOnBlur, validateOnChange } = fieldContext

    const hasError = $derived(form.hasFieldErrors(name as never))

    // Generate IDs for accessibility
    const inputId = `field-${name}`
    const errorId = `field-${name}-error`
    const descriptionId = `field-${name}-description`

    // Event handlers
    function handleBlur(event: FocusEvent) {
        if (validateOnBlur) {
            form.validateField(name as never)
        }
    }

    function handleInput(event: Event) {
        // Always clear errors when user types
        if (form.state.errors[name]) {
            delete form.state.errors[name]
            // Recalculate valid state
            form.state.valid = Object.keys(form.state.errors).length === 0
        }
        
        // Optionally validate on change
        if (validateOnChange) {
            form.validateField(name as never)
        }
    }

    const props = $derived({
        id: inputId,
        name: name,
        "aria-invalid": hasError,
        "aria-describedby": `${descriptionId} ${hasError ? errorId : ""}`.trim(),
        onblur: handleBlur,
        oninput: handleInput,
    })
</script>

<div class={cls}>
    {@render children({ props })}
</div>
