export type ServerContext<T> = {
    id: string
    ids: Record<string, string>
    data: T
}