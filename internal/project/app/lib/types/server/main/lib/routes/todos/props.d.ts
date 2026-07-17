export type Props = todos.Props

export declare namespace todos {
    export type Props = {
        items: null|(schema.Todo[])
        error: string
    }
}

export declare namespace schema {
    export type Todo = {
        id: string
        sessionId: string
        description: string
        checked: number
    }
}