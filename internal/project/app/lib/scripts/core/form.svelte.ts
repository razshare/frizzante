import type { form } from "$gen/types/main/lib/core/form/State"

/**
 * Form store state that mirrors the Go form.State type
 */
export type FormState<T = Record<string, unknown>> = {
    data: T
    errors: Record<string, string[]>
    valid: boolean
    tainted: Record<string, boolean>
    message: string
    pending: boolean
}

/**
 * Validation rule function
 */
export type ValidationRule<T = unknown> = (value: T) => string | null

/**
 * Validation schema - maps field names to validation rules
 */
export type ValidationSchema<T = Record<string, unknown>> = {
    [K in keyof T]?: ValidationRule<T[K]>[]
}

/**
 * Form context for nested components
 */
export class FormContext<T = Record<string, unknown>> {
    state: FormState<T> = $state({
        data: {} as T,
        errors: {},
        valid: true,
        tainted: {},
        message: "",
        pending: false,
    })

    schema: ValidationSchema<T> = {}

    constructor(initialData: T, schema?: ValidationSchema<T>) {
        this.state.data = initialData
        this.schema = schema || {}
    }

    /**
     * Set a field value and mark it as tainted
     */
    setField<K extends keyof T>(name: K, value: T[K]): void {
        this.state.data[name] = value
        this.state.tainted[name as string] = true

        // Clear errors for this field when user starts typing
        if (this.state.errors[name as string]) {
            delete this.state.errors[name as string]
            
            // Recalculate valid state
            this.state.valid = Object.keys(this.state.errors).length === 0
        }    }

    /**
     * Validate a single field
     */
    validateField<K extends keyof T>(name: K): boolean {
        const rules = this.schema[name]
        if (!rules) return true

        const value = this.state.data[name]
        const errors: string[] = []

        for (const rule of rules) {
            const error = rule(value)
            if (error) {
                errors.push(error)
            }
        }

        if (errors.length > 0) {
            this.state.errors[name as string] = errors
            return false
        } else {
            delete this.state.errors[name as string]
            return true
        }
    }

    /**
     * Validate all fields
     */
    validate(): boolean {
        let isValid = true

        for (const fieldName in this.schema) {
            if (!this.validateField(fieldName)) {
                isValid = false
            }
        }

        this.state.valid = isValid
        return isValid
    }

    /**
     * Reset the form to initial state
     */
    reset(initialData: T): void {
        this.state.data = initialData
        this.state.errors = {}
        this.state.valid = true
        this.state.tainted = {}
        this.state.message = ""
        this.state.pending = false
    }

    /**
     * Update form state from server response
     */
    updateFromServer(serverState: Partial<FormState<T>>): void {
        if (serverState.data !== undefined) this.state.data = serverState.data
        if (serverState.errors !== undefined) this.state.errors = serverState.errors
        if (serverState.valid !== undefined) this.state.valid = serverState.valid
        if (serverState.tainted !== undefined) this.state.tainted = serverState.tainted
        if (serverState.message !== undefined) this.state.message = serverState.message
        if (serverState.pending !== undefined) this.state.pending = serverState.pending
    }

    /**
     * Check if any field has been tainted
     */
    get isTainted(): boolean {
        return Object.keys(this.state.tainted).length > 0
    }

    /**
     * Get errors for a specific field
     */
    getFieldErrors(name: keyof T): string[] {
        return this.state.errors[name as string] || []
    }

    /**
     * Check if a field has errors
     */
    hasFieldErrors(name: keyof T): boolean {
        const errors = this.state.errors[name as string]
        return errors !== undefined && errors.length > 0
    }

    /**
     * Check if a field is tainted
     */
    isFieldTainted(name: keyof T): boolean {
        return this.state.tainted[name as string] === true
    }
}

/**
 * Create a new form store
 *
 * Example:
 * ```ts
 * const form = superForm({ email: '', password: '' }, {
 *   email: [required(), email()],
 *   password: [required(), minLength(8)]
 * })
 * ```
 */
export function superForm<T = Record<string, unknown>>(
    initialData: T,
    schema?: ValidationSchema<T>
): FormContext<T> {
    return new FormContext(initialData, schema)
}

/**
 * Built-in validation rules
 */

export function required(message = "This field is required"): ValidationRule {
    return (value: unknown): string | null => {
        if (value === null || value === undefined || value === "") {
            return message
        }
        return null
    }
}

export function email(message = "Invalid email address"): ValidationRule<string> {
    return (value: string): string | null => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        if (value && !emailRegex.test(value)) {
            return message
        }
        return null
    }
}

export function minLength(min: number, message?: string): ValidationRule<string> {
    return (value: string): string | null => {
        if (value && value.length < min) {
            return message || `Must be at least ${min} characters`
        }
        return null
    }
}

export function maxLength(max: number, message?: string): ValidationRule<string> {
    return (value: string): string | null => {
        if (value && value.length > max) {
            return message || `Must be at most ${max} characters`
        }
        return null
    }
}

export function min(minValue: number, message?: string): ValidationRule<number> {
    return (value: number): string | null => {
        if (value < minValue) {
            return message || `Must be at least ${minValue}`
        }
        return null
    }
}

export function max(maxValue: number, message?: string): ValidationRule<number> {
    return (value: number): string | null => {
        if (value > maxValue) {
            return message || `Must be at most ${maxValue}`
        }
        return null
    }
}

export function pattern(regex: RegExp, message = "Invalid format"): ValidationRule<string> {
    return (value: string): string | null => {
        if (value && !regex.test(value)) {
            return message
        }
        return null
    }
}

export function url(message = "Invalid URL"): ValidationRule<string> {
    return (value: string): string | null => {
        try {
            new URL(value)
            return null
        } catch {
            return message
        }
    }
}
