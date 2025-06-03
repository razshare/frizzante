import { hydrate } from "svelte";
import RenderClient from "./client.svelte";
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
target().innerHTML = "";
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
hydrate(RenderClient, { target: target(), props: props() });
