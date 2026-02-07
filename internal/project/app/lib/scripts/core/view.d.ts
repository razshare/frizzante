export type View<T> = {
    props: T
    name: string
    align: number
    render: number
    pin: () => Promise<void>
    type: "" | "default" | "snapshot"
}
