package runtime

func init() {
	makeAnyIns()

	makeNullIns()

	makeNumeric()
	makeIntIns()

	makeStringIns()

	makeGoFuncIns()

	makeGorillaToInteger()
	makeGorillaToString()

	makeGlobal()
}
