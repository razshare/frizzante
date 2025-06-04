import Welcome from './Welcome.svelte'
import type {Component} from "svelte";

export const views: Record<string, Component> = {
    "Welcome": Welcome,
}