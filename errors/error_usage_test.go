package errors

import (
	"fmt"
	"testing"
)

var myerror = New("name is empty")

func Service(name string) error {
	if name == "" {
		return myerror
	}
	return nil
}

func Controller() error {
	err := Service("")
	if err != nil {
		return Wrap(err, "Service has error")
	}
	return nil
}

func funcA() error {
	if err := funcB(); err != nil {
		return Wrap(err, "call funcB failed")
	}
	return New("func called error")
}

func funcB() error {
	return New("func called error")
}

// 自定义的错误类型
type DefineError struct {
	msg string
}

func (d *DefineError) Error() string {
	return d.msg
}

func Service1(name string) error {
	if name == "" {
		return &DefineError{msg: "error 1"}
	}
	return nil
}

func Controller1() error {
	err := Service1("")
	if err != nil {
		return Wrap(err, "Service has error")
	}
	return nil
}

/*
测试Wrap 函数

=== RUN   TestWrap

	error_user_test.go:21: call funcB failed

=============

	error_user_test.go:25: func called error
	    logtest.funcB
	    	/Users/maozhongyu/code/logtest/error_user_test.go:17
	    logtest.funcA
	    	/Users/maozhongyu/code/logtest/error_user_test.go:10
	    logtest.TestWrap
	    	/Users/maozhongyu/code/logtest/error_user_test.go:25
	    testing.tRunner
	    	/Users/maozhongyu/go1.18.5/src/testing/testing.go:1439
	    runtime.goexit
	    	/Users/maozhongyu/go1.18.5/src/runtime/asm_amd64.s:1571
	    call funcB failed
	    logtest.funcA
	    	/Users/maozhongyu/code/logtest/error_user_test.go:11
	    logtest.TestWrap
	    	/Users/maozhongyu/code/logtest/error_user_test.go:25
	    testing.tRunner
	    	/Users/maozhongyu/go1.18.5/src/testing/testing.go:1439
	    runtime.goexit
	    	/Users/maozhongyu/go1.18.5/src/runtime/asm_amd64.s:1571

--- PASS: TestWrap (0.00s)
PASS
*/
func TestWrap(t *testing.T) {
	t.Logf("%v", funcA())

	fmt.Println("=============")

	t.Logf("%+v", funcA())

}

/*
测试 Is函数

一般用在哨兵错误

=== RUN   TestIs
is myerror value
--- PASS: TestIs (0.00s)
PASS
*/
func TestIs(t *testing.T) {
	err := Controller()
	if Is(err, myerror) {
		fmt.Println("is myerror value") // 输出
	}
}

/*
测试 As函数
error 1
Service has error
logtest.Controller1

	/Users/maozhongyu/code/logtest/error_user_test.go:56

logtest.TestAs

	/Users/maozhongyu/code/logtest/error_user_test.go:116

testing.tRunner

	/Users/maozhongyu/go1.18.5/src/testing/testing.go:1439

runtime.goexit

	/Users/maozhongyu/go1.18.5/src/runtime/asm_amd64.s:1571

-----------
AS error
*/
func TestAs(t *testing.T) {
	err := Controller1()
	fmt.Printf("%+v\n", err)
	fmt.Println("-----------")
	var myerror *DefineError
	if As(err, &myerror) {
		fmt.Println("AS error")
	}
}
