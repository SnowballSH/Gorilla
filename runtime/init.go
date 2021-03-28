package runtime

func init() {
	makeIntIns()
	makeStringIns()

	makeGorillaToInteger()
	makeGorillaToString()

	makeGlobal()
}
