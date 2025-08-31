import { type Readable, readable } from "svelte/store"
import { IS_BROWSER } from "$lib/scripts/core/constants.ts"

export function source(path: string) {
    if (!IS_BROWSER) {
        return {
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            select(event: string = "message"): Readable<string> {
                return readable("")
            },
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            selectJson<T>(event: string = "message"): Readable<undefined | T> {
                return readable(undefined)
            },
        }
    }

    const source = new EventSource(path)

    return {
        select(event: string = "message"): Readable<string> {
            return readable<string>("", set => {
                source.addEventListener(event, ev => {
                    set(ev.data)
                })
                return function stop() {
                    source.close()
                }
            })
        },
        selectJson<T>(event: string = "message"): Readable<T> {
            return readable<T>(undefined, set => {
                source.addEventListener(event, ev => {
                    set(JSON.parse(ev.data))
                })
                return function stop() {
                    source.close()
                }
            })
        },
    }
}
