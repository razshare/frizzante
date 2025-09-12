import { mount } from "svelte"
import App from "./app.client.svelte"
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
target().innerHTML = ""
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
mount(App, { target: target(), props: data() })
