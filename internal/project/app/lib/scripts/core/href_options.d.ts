export type HrefOptions = {
    onpending?: () => void
    ondone?: () => void
    onerror?: (error: Error) => void
    pendingDelay?: number
}
