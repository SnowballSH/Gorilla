package runtime

func init() {
	makeNullIns()
	makeIntIns()
	makeStringIns()

	makeGorillaToInteger()
	makeGorillaToString()

	makeGlobal()
}
