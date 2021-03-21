package objects

type Environment struct {
	store                map[string]Object
	externalEnvironments map[string]*Environment
	outer                *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store:                map[string]Object{},
		externalEnvironments: map[string]*Environment{},
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	env.externalEnvironments = outer.externalEnvironments
	return env
}

func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e.store[key]
	if ok {
		return obj, ok
	}

	if e.outer == nil {
		return nil, false
	}

	return e.outer.Get(key)
}

func (e *Environment) Set(key string, value Object) Object {
	e.store[key] = value
	return value
}

func (e *Environment) AddExternalEnvironment(alias string, env *Environment) {
	e.externalEnvironments[alias] = env
}

func (e *Environment) GetExternalEnvironment(alias string) (*Environment, bool) {
	env, ok := e.externalEnvironments[alias]
	return env, ok
}
