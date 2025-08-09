export type View<T> = {
    name: string
    data: T
    renderMode: number
}

export type SwapConfig = {
    method: string
    path: string
    body: false | FormData
    view: View<unknown>
}
