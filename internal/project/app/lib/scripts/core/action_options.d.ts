export type ActionOptions = {
    onpending?: () => void
    ondone?: () => void
    onerror?: (error: Error) => void
    pendingDelay?: number
}
