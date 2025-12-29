package actions

type MigrateOptions struct {
	Strict   bool
	SqlcYaml string
	Query    string
	Database string
}
