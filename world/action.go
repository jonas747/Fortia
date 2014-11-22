package world

type Action struct {
	Name       string
	Handler    string
	Properties map[string]interface{}
}

func (a *Action) GetPropertyInt(prop string) int {
	val := a.Properties[prop].(int)
	return val
}

func (a *Action) GetPropertyString(prop string) string {
	val := a.Properties[prop].(string)
	return val
}
