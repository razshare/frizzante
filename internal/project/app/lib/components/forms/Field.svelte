<script lang="ts">
    import type { Snippet } from "svelte"
    import { getContext, setContext } from "svelte"
    import type { FormContext } from "$lib/scripts/core/form.svelte.ts"

    type Props<T = Record<string, unknown>> = {
        form?: FormContext<T>
        name: string
        children: Snippet<
            [
                {
                    errors: string[]
                    tainted: boolean
                    hasError: boolean
                },
            ]
        >
        class?: string
        validateOnBlur?: boolean
        validateOnChange?: boolean
    }

    let {
        form: formProp,
        name,
        children,
        class: cls,
        validateOnBlur = true,
        validateOnChange = false,
    }: Props = $props()

    // Get form from context if not provided as prop
    const formContext = formProp || (getContext("superform") as FormContext)

    if (!formContext) {
        throw new Error("Field component must be used within a SuperForm or have form prop")
    }

    // Set field context for nested components (Label, Control, FieldErrors, etc.)
    setContext("field", {
        name,
        form: formContext,
        validateOnBlur,
        validateOnChange,
    })

    // Reactive field state
    const errors = $derived(formContext.getFieldErrors(name as never))
    const tainted = $derived(formContext.isFieldTainted(name as never))
    const hasError = $derived(formContext.hasFieldErrors(name as never))
</script>

<div class={cls} data-field={name} data-has-error={hasError} data-tainted={tainted}>
    {@render children({ errors, tainted, hasError })}
</div>
