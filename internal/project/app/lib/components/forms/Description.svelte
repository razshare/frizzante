<script lang="ts">
    import type { Snippet } from "svelte"
    import { getContext } from "svelte"

    type Props = {
        children: Snippet
        class?: string
    }

    let { children, class: cls }: Props = $props()

    // Get field context
    const fieldContext = getContext("field") as { name: string } | undefined

    const descriptionId = fieldContext ? `field-${fieldContext.name}-description` : undefined
</script>

<div id={descriptionId} class={cls || "text-sm text-base-content/60 mt-1"}>
    {@render children()}
</div>
