package junklang

type Scope struct {
    variables map[string]interface{}
    functions map[string]interface{}
    parent    *Scope
}

func NewScope(parent *Scope) *Scope {
    return &Scope{
        variables: make(map[string]interface{}),
        functions: make(map[string]interface{}),
        parent:    parent,
    }
}

func (s *Scope) Set(name string, value interface{}) {
    s.variables[name] = value
}

func (s *Scope) Get(name string) interface{} {
    if val, ok := s.variables[name]; ok {
        return val
    }
    if s.parent != nil {
        return s.parent.Get(name)
    }
    panic("udefineret variabel: " + name)
}

func (s *Scope) SetFunction(name string, fn interface{}) {
    s.functions[name] = fn
}

func (s *Scope) GetFunction(name string) interface{} {
    if fn, ok := s.functions[name]; ok {
        return fn
    }
    if s.parent != nil {
        return s.parent.GetFunction(name)
    }
    panic("udefineret funktion: " + name)
}