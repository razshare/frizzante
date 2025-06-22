import { hydrate } from "svelte"
import ClientView from "$lib/utilities/frz/components/ClientView.svelte"
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
target().innerHTML = ""
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
hydrate(ClientView, { target: target(), props: props() })
