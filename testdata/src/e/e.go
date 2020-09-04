package e

type T struct {
	Hoge string
}

/* -----UseExpectedを生成する----- */

// offset_assign
func (t *T) assgin() {
	t.Hoge = "hoge"
}

/* -----UseExpectedを生成しない----- */

// offset_sameTypeDiffVar
func (t *T) sameTypeDiffVar() {
	h := &T{}
	h.Hoge = "hoge"
}

// offset_assignInMethod
func (t *T) assignInMethod() {
	t.assgin()
}

func setHoge(t *T) {
	t.Hoge = "hoge"
}

// offset_assignInFunc
func (t *T) assignInFunc() {
	setHoge(t)
}

// offset_assignInGoFunc
func (t *T) assignInGoFunc() {
	go func() {
		t.Hoge = "hoge"
	}()
}
