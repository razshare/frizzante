import { hydrate } from "svelte";
import RenderClient from "./client.svelte";
// @ts-ignore
target().innerHTML = "";
// @ts-ignore
hydrate(RenderClient, { target: target(), props: props() });
