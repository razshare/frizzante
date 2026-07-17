import type { Root } from "$lib/types/core/root"

export const root = $state<Root>({
    view: {
        name: "",
        ondone() {
            return {}
        },
    },
    type: "",
    swapping: false,
})
