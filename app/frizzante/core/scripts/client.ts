import { mount } from "svelte"
import ClientView from "$frizzante/core/components/ClientView.svelte"
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
target().innerHTML = ""
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
mount(ClientView, { target: target(), props: props() })
