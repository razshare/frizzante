<script lang="ts">
    import type { Snippet } from "svelte"
    import { getContext } from "svelte"

    type Props = {
        children: Snippet
        class?: string
        for?: string
    }

    let { children, class: cls, for: forProp }: Props = $props()

    // Get field context
    const fieldContext = getContext("field") as { name: string } | undefined

    // Use explicit 'for' prop if provided, otherwise derive from field context
    const forId = $derived(forProp || (fieldContext ? `field-${fieldContext.name}` : undefined))
</script>

<label for={forId} class={cls}>
    {@render children()}
</label>
