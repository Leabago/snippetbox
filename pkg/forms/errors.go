package forms

type errors map[string][]string

func (e errors) Add(filed string, message string) {
	e[filed] = append(e[filed], message)
}

func (e errors) Get(filed string) string {
	es := e[filed]
	if len(es) == 0 {
		return ""
	} else {
		return es[0]
	}
}
