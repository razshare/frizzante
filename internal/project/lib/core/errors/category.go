package errors

type Category string

const (
	CategoryValidation Category = "validation"
	CategoryNotFound Category = "not_found"
	CategoryPermission Category = "permission"
	CategoryConflict Category = "conflict"
	CategoryInput Category = "input"
	CategoryInternal Category = "internal"
	CategoryTimeout Category = "timeout"
	CategoryRateLimit Category = "rate_limit"
	CategoryConfiguration Category = "configuration"
	CategoryDependency Category = "dependency"
)