import { render as _render } from "svelte/server";
import RenderServer from "./server.svelte";
export async function render(props:any) {
  return _render(RenderServer, { props });
}
