export type View<T> = {
    name: string
    data: T
    renderMode: number
}

export type HistoryEntry = {
    nodeName: string
    method: string
    url: string
    body: Record<string, string>
}
