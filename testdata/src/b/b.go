package b

// offset_intList
func intList(list []int) []int {
	return list
}

// offset_mapFunc
func mapFunc(input map[int]string) (map[int]string, map[string]error) {
	return input, map[string]error{}
}

// offset_pointer
func pointer(input *string) *string {
	return input
}

// offset_pointerList
func pointerList(input []*string) []*string {
	return input
}

// offset_function
func function(input func(i int) string) func(i int) string {
	return input
}

// offset_chanel
func chanel(input chan int) chan int {
	return input
}

type myStruct struct{}

// offset_myStructFunc
func myStructFunc(ms myStruct) myStruct {
	return ms
}
