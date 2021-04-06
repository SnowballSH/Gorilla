package runtime

var NumericClass *RClass

func makeNumeric() {
	NumericClass = MakeClassFromSuper("Numeric", AnyClass, NotCallable, NewEnvironment())
}
