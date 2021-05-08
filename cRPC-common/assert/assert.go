package assert

func CheckParam(expr bool) {
	Assert(expr, "invalid param")
}

func Assert(expr bool, message string) {
	if !expr {
		panic(message)
	}
}