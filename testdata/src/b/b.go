package b

func intList(list []int) []int {
	return list
}

func mapFunc(input map[int]string) (map[int]string, map[string]error) {
	return input, map[string]error{}
}

func pointer(input *string) *string {
	return input
}

func pointerList(input []*string) []*string {
	return input
}

func function(input func(i int) string) func(i int) string {
	return input
}

func chanel(input chan int) chan int {
	return input
}
