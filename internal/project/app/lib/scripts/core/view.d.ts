import { views as clientViews } from "$exports.client"
import { views as serverViews } from "$exports.server"

export type View = {
    align: number
    render: number
    props: Record<string, unknown>
    type: "" | "default" | "snapshot"
    name: keyof typeof clientViews | keyof typeof serverViews
}
