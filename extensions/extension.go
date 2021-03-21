package extensions

type Plugin interface {
	Eval(fnName string, args ...interface{}) ([]interface{}, error)
	Package() string
	Version() string
	Description() string
}
