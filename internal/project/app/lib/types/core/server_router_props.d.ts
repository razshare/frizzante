import type { views } from "$exports.server"

export type ServerRouterProps = {
    name: keyof typeof views
    props: Record<string, unknown>
    type: "" | "default" | "snapshot"
}
