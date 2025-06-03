import { render as _render } from "svelte/server";
import RenderServer from "./server.svelte";
export async function render(props:unknown) {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-expect-error
  return _render(RenderServer, { props });
}
