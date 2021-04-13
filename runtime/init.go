package runtime

func init() {
	makeAnyIns()

	makeNullIns()
	makeBoolIns()

	makeNumeric()
	makeIntIns()

	makeStringIns()

	makeGoFuncIns()
	makeLambdaIns()

	makeGorillaToInteger()
	makeGorillaToString()

	makeGlobal()
}
