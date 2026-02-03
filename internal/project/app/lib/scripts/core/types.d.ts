export type View<T> = {
    props: T
    name: string
    align: number
    render: number
    pin: () => Promise<void>
}

export type HistoryEntry = {
    nodeName: string
    method: string
    url: string
    body: Record<string, string>
}
