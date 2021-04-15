package runtime

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/grammar"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLambda(t *testing.T) {
	w := NewVM(nil)

	x := []byte{
		grammar.Magic,
		grammar.GetVar, 1, 'x',

		grammar.Integer, 1, 0x02, // 2

		grammar.GetInstance, // 2.*
		1, '*',              // *

		grammar.Call, // 2.*(3)
		1, 0x01,      // 1 arg

		grammar.Pop,
	}
	lambda := NewLambda([]string{"x"}, x, w)
	assert.EqualValues(t, x, lambda.InternalValue)
	assert.EqualValues(t, x, lambda.Value())
	assert.Equal(t, "Lambda Function", lambda.ToString())
	assert.Equal(t, fmt.Sprintf("Lambda Function %p", lambda), lambda.Inspect())
	assert.Equal(t, LambdaClass, lambda.Class())
	assert.True(t, lambda.IsTruthy())
	assert.False(t, lambda.EqualTo(NewLambda(nil, nil, w)))
	assert.False(t, lambda.EqualTo(nil))

	wot := NewString("")
	assert.False(t, lambda.EqualTo(wot))

	k, e := lambda.Call(lambda, NewInteger(3))
	assert.Nil(t, e)
	assert.Equal(t, "6", k.ToString())

	x = []byte{
		grammar.Magic,
	}
	lambda = NewLambda(nil, x, w)

	k, e = lambda.Call(lambda)
	assert.Nil(t, e)
	assert.Equal(t, "null", k.ToString())

	x = []byte{
		grammar.Magic,
		grammar.GetVar, 1, '%',

		grammar.Integer, 1, 0x02,

		grammar.GetInstance,
		1, '*',

		grammar.Call,
		1, 0x01,

		grammar.Pop,
	}
	lambda = NewLambda(nil, x, w)
	k, e = lambda.Call(lambda)
	assert.NotNil(t, e)

	x = []byte{
		grammar.Magic,
	}
	lambda = NewLambda([]string{"x"}, x, w)
	k, e = lambda.Call(lambda)
	assert.NotNil(t, e)
}
