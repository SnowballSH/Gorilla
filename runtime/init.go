package runtime

func init() {
	makeAnyIns()

	makeNullIns()
	makeBoolIns()

	makeNumeric()
	makeIntIns()

	makeStringIns()

	makeGoFuncIns()

	makeGorillaToInteger()
	makeGorillaToString()

	makeGlobal()
}
