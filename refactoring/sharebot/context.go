package sharebot

// Context container witch contains results needed for bot work
type Context struct {
	TGToken string
}

// GetContext returns actual contecxt
func GetContext() (Context, error) {
	panic("implement me")
}
