type ServerProperties<T> = {
    id: string
    ids: Record<string, string>
    data: T
}