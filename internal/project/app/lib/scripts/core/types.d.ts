export type View<T> = {
    name: string
    props: T
    render: number
    align: number
    pending: boolean
}

export type HistoryEntry = {
    nodeName: string
    method: string
    url: string
    body: Record<string, string>
}
