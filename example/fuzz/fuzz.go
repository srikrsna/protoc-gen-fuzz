package fuzz

var ff = [...]interface{}{
	FuzzSomeMessage,
	FuzzSomeMessage_OneofString,
	FuzzSomeMessage_OneofBool,
	FuzzInnerMessage,
}

func FuzzFuncs() []interface{} {
	return ff[:]
}
