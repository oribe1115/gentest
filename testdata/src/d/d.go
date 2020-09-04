package d

type T struct{}

// offset_basicRecv
func (t T) basicRecv() {}

// offset_pointerRecv
func (t *T) pointerRecv() {}

// offset_paralell
func parallel() {}
