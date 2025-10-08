<script lang="ts">
    import type { Snippet } from "svelte"
    import { getContext, setContext } from "svelte"
    import { swap } from "$lib/scripts/core/swap.ts"
    import { route } from "$lib/scripts/core/route.ts"
    import type { View } from "$lib/scripts/core/types"
    import type { FormContext } from "$lib/scripts/core/form.svelte.ts"
    import { IS_BROWSER } from "$lib/scripts/core/constants.ts"

    type Props<T = Record<string, unknown>> = {
        form: FormContext<T>
        action?: string
        method?: "GET" | "POST"
        enctype?:
            | "multipart/form-data"
            | "application/x-www-form-urlencoded"
            | "text/plain"
        children: Snippet<
            [
                {
                    pending: boolean
                    valid: boolean
                    errors: Record<string, string[]>
                    message: string
                    tainted: boolean
                },
            ]
        >
        class?: string
        style?: string
        onSubmit?: (event: SubmitEvent) => void | Promise<void>
        validateOnSubmit?: boolean
    }

    let {
        form,
        action: actionPath = "",
        method = "POST",
        enctype,
        children,
        class: cls,
        style,
        onSubmit,
        validateOnSubmit = true,
    }: Props = $props()

    setContext("superform", form)
    
    // Get view context during component initialization
    const view = IS_BROWSER ? (getContext("view") as View<never>) : null

    async function handleSubmit(event: SubmitEvent) {
        event.preventDefault()

        if (validateOnSubmit && !form.validate()) {
            return
        }

        form.state.pending = true

        try {
            if (onSubmit) {
                await onSubmit(event)
            }

            if (IS_BROWSER && view) {
                route(view)
                const formElement = event.target as HTMLFormElement
                const record = await swap(formElement, view)
                record()
            }
        } catch (error) {
            console.error("Form submission error:", error)
            form.state.message = "An error occurred during submission"
            form.state.valid = false
        } finally {
            form.state.pending = false
        }
    }
</script>

<form {enctype} {method} action={actionPath} onsubmit={handleSubmit} class={cls} {style}>
    {@render children({
        pending: form.state.pending,
        valid: form.state.valid,
        errors: form.state.errors,
        message: form.state.message,
        tainted: form.isTainted,
    })}
</form>
