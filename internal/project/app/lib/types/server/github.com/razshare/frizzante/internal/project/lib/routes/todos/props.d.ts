export type Props = todos.Props

export declare namespace schema {
    export type Todo = {
        id: string
        sessionId: string
        description: string
        checked: number
    }
}

export declare namespace todos {
    export type Props = {
        items: null|(schema.Todo[])
        error: string
    }
}