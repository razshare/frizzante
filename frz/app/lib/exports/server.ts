import Welcome from "$lib/components/views/Welcome.svelte"
import type { Component } from "svelte"

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
export const views: Record<string, Component<unknown>> = {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    "Welcome": Welcome,
}
