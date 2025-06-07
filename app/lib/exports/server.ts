import Welcome from '$lib/components/views/Welcome.svelte'
import type {Component} from "svelte";

export const views: Record<string, Component> = {
    "Welcome": Welcome,
}