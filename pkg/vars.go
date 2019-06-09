package pkg

// VarStorer is an abstraction that is used to set or get variables and it's values.
type VarStorer interface {
	Register(name, val string) error
	Replace(value string) string
	Var(name string) (string, error)
}
