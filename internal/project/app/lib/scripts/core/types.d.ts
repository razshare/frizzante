export type View<T> = {
    Name: string
    Props: T
    Render: number
    Align: number
}

export type HistoryEntry = {
    nodeName: string
    method: string
    url: string
    body: Record<string, string>
}
