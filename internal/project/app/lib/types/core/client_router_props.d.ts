import type { views } from "$exports.client"

export type ClientRouterProps = {
    name: keyof typeof views
    props: Record<string, unknown>
    type: "" | "default" | "snapshot"
}
