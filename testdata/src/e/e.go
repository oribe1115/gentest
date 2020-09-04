package e

type T struct {
	Hoge string
}

type P struct {
	Po string
}

// offset_assign
func (t *T) assgin() {
	t.Hoge = "hoge"
}

// offset_sameTypeDiffVar
func (t *T) sameTypeDiffVar() {
	h := &T{}
	h.Hoge = "hoge"
}
