import type { views as clientViews } from "$exports.client"
import type { views as serverViews } from "$exports.server"

export type Root = {
    view: {
        ondone: () => Record<string, unknown>
        name: keyof typeof clientViews | keyof typeof serverViews | ""
    }
    type: "" | "default" | "snapshot"
    swapping: boolean
}
