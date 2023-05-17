## errors包

基于 `github.com/pkg/errors` 包，增加对 `error code` 的支持，完全兼容 `github.com/pkg/errors`

并增加grpc error 和 with code error 的互转


### Usage



**Wrap,Is,As函数使用**
```go

import (
	"fmt"
	"testing"
	
	"github.com/cr-mao/lori/errors"
)

var myerror = errors.New("name is empty")

func Service(name string) error {
	if name == "" {
		return myerror
	}
	return nil
}

func Controller() error {
	err := Service("")
	if err != nil {
		return errors.Wrap(err, "Service has error")
	}
	return nil
}

func funcA() error {
	if err := funcB(); err != nil {
		return errors.Wrap(err, "call funcB failed")
	}
	return errors.New("func called error")
}

func funcB() error {
	return errors.New("func called error")
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
		return errors.Wrap(err, "Service has error")
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
	if errors.Is(err, myerror) {
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
	if errors.As(err, &myerror) {
		fmt.Println("AS error")
	}
}
```

### WithCode,WrapC,IsCode,Cause 使用

```go

import (
	"fmt"
	"testing"

	"github.com/cr-mao/lori/errors"
)

const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 110001

	// ErrUserAlreadyExist - 400: User encode error.
	ErrUserEncode
)

// ErrCode implements `github.com/cr-mao/errors`.Coder interface.
type ErrCode struct {
	// C refers to the code of the ErrCode.
	C int

	// HTTP status that should be used for the associated error code.
	HTTP int

	// External (user) facing error text.
	Ext string

	// Ref specify the reference document.
	Ref string
}

// Code returns the integer code of ErrCode.
func (coder ErrCode) Code() int {
	return coder.C
}

// String implements stringer. String returns the external error message,
// if any.
func (coder ErrCode) String() string {
	return coder.Ext
}

// Reference returns the reference document.
func (coder ErrCode) Reference() string {
	return coder.Ref
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder ErrCode) HTTPStatus() int {
	if coder.HTTP == 0 {
		return 500
	}
	return coder.HTTP
}

func register(code int, httpStatus int, message string, refs ...string) {
	if httpStatus != 200 && httpStatus != 400 && httpStatus != 401 && httpStatus != 403 && httpStatus != 404 && httpStatus != 500 {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}
	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	}
	errors.MustRegister(coder)
}

func TestWithCode(t *testing.T) {

	//注册code
	register(ErrUserNotFound, 400, "user not found")
	register(ErrUserEncode, 400, "user  encode error")

	err := errors.WithCode(ErrUserNotFound, "user %s not found.", "cr-mao")
	if err != nil {
		fmt.Println(err.Error()) // user not found
	}
}

func TestWithCodeDetail(t *testing.T) {
	//注册code
	register(ErrUserNotFound, 400, "user not found")
	register(ErrUserEncode, 400, "user encode error")
	if err := bindUser(); err != nil {
		// %s: Returns the user-safe error string mapped to the error code or the error message if none is specified.
		fmt.Println("====================> s <====================")
		fmt.Printf("%s\n\n", err)

		// %v: Alias for %s.
		fmt.Println("====================> v <====================")
		fmt.Printf("%v\n\n", err)

		// %-v: Output caller details, useful for troubleshooting.
		fmt.Println("====================> -v <====================")
		fmt.Printf("%-v\n\n", err)

		// %+v: Output full error stack details, useful for debugging.
		fmt.Println("====================> +v <====================")
		fmt.Printf("%+v\n\n", err)

		// %#-v: Output caller details, useful for troubleshooting with JSON formatted output.
		fmt.Println("====================> #-v <====================")
		fmt.Printf("%#-v\n\n", err)

		// %#+v: Output full error stack details, useful for debugging with JSON formatted output.
		fmt.Println("====================> #+v <====================")
		fmt.Printf("%#+v\n\n", err)

		// do some business process based on the error type
		if errors.IsCode(err, ErrUserEncode) {
			fmt.Println("this is a ErrEncodingFailed error") // this is a ErrEncodingFailed error
		}

		if errors.IsCode(err, ErrUserNotFound) {
			fmt.Println("this is a ErrDatabase error") // this is a ErrDatabase error
		}

		// we can also find the cause error
		fmt.Println(errors.Cause(err)) // user not found
	}
}

func bindUser() error {
	if err := getUser(); err != nil {
		// Step3: Wrap the error with a new error message and a new error code if needed.
		return errors.WrapC(err, ErrUserEncode, "encoding user %s failed.", "cr-mao")
	}

	return nil
}

func getUser() error {
	if err := queryDatabase(); err != nil {
		// Step2: Wrap the error with a new error message.
		return errors.Wrap(err, "get user failed.")
	}
	return nil
}

func queryDatabase() error {
	// Step1. Create error with specified error code.
	return errors.WithCode(ErrUserNotFound, "user %s not found.", "cr-mao")
}

```

结果：
```text
====================> s <====================
user encode error

====================> v <====================
user encode error

====================> -v <====================
encoding user cr-mao failed. - #2 [/Users/maozhongyu/code/game-server/tmp_test.go:134 (game-server.bindUser)] (110002) user encode error

====================> +v <====================
encoding user cr-mao failed. - #2 [/Users/maozhongyu/code/game-server/tmp_test.go:134 (game-server.bindUser)] (110002) user encode error; get user failed. - #1 [/Users/maozhongyu/code/game-server/tmp_test.go:143 (game-server.getUser)] (110001) user not found; user cr-mao not found. - #0 [/Users/maozhongyu/code/game-server/tmp_test.go:150 (game-server.queryDatabase)] (110001) user not found

====================> #-v <====================
[{"caller":"#2 /Users/maozhongyu/code/game-server/tmp_test.go:134 (game-server.bindUser)","code":110002,"error":"encoding user cr-mao failed.","message":"user encode error"}]

====================> #+v <====================
[{"caller":"#2 /Users/maozhongyu/code/game-server/tmp_test.go:134 (game-server.bindUser)","code":110002,"error":"encoding user cr-mao failed.","message":"user encode error"},{"caller":"#1 /Users/maozhongyu/code/game-server/tmp_test.go:143 (game-server.getUser)","code":110001,"error":"get user failed.","message":"user not found"},{"caller":"#0 /Users/maozhongyu/code/game-server/tmp_test.go:150 (game-server.queryDatabase)","code":110001,"error":"user cr-mao not found.","message":"user not found"}]

this is a ErrEncodingFailed error
this is a ErrDatabase error
user not found
```



### WithCode格式占位符

withCode错误实现了一个`func (w *withCode) Format(state fmt.State, verb rune)`方法，该方法用来打印不同格式的错误信息，见下表：

| 格式占位符 | 格式化操作                                                   |
| ---------- | ------------------------------------------------------------ |
| %s         | 返回可以直接展示给用户的错误信息                             |
| %v         | alias for %s                                                 |
| %-v        | 打印出调用栈，错误码、展示给用户的错误信息、展示给研发的错误信息（只展示错误链中的最有一个错误） |
| %+v        | 打印出调用栈，错误码、展示给用户的错误信息、展示给研发的错误信息（展示错误链中所有错误） |
| %#-v       | json格式打印出调用栈，错误码、展示给用户的错误信息、展示给研发的错误信息（只展示错误链中的最有一个错误） |
| %#+v       | json格式打印出调用栈，错误码、展示给用户的错误信息、展示给研发的错误信息（展示错误链中所有错误） |



### 错误码设计

项目组代号:10

服务代号:01

模块代号:0~99

错误码：0~99

| 错误标识                | 错误码   | HTTP状态码 | 描述                          |
| ----------------------- | -------- | ---------- | ----------------------------- |
| ErrNo                   | 10010000 | 200        |  OK                            |
| ErrInternalServer       | 10010001 | 500        |  Internal server error （服务器内部错误）      |
| ErrParams               | 10010002 | 400        |  Illegal params  (请求参数不合法)                |
| ErrAuthenticationHeader | 10010003 | 401        |  Authentication header Illegal  (要登录的接口，头的token认证失败,失败跳登录页面)|
| ErrAuthentication       | 10010004 | 401        |  Authentication failed  (登录失败，输入账户、密码失败)|
| ErrNotFound             | 10010005 | 404        |  Route not found     (请求路由找不到）             |
| ErrPermission           | 10010006 | 403        |  Permission denied (没有权限,一些接口可能没请求权限， 这个估计暂时用不到)            |
| ErrTooFast              | 10010007 | 429        |  Too Many Requests （用户在给定的时间内发送了太多请求）            |
| ErrTimeout              | 10010008 | 504        |  Server response timeout   （go服务这边不会返回，一般是nginx、网关超时 才返回504）|
| ErrMysqlServer          | 10010101 | 500        |  Mysql server error      （mysql 服务错误)       |
| ErrMysqlSQL             | 10010102 | 500        |  Illegal SQL               (sql 代码错误）       |
| ErrRedisServer          | 10010201 | 500        |  Redis server error        （redis 服务错误）    |


http code标志错误，非200， 客户端再根据 错误码进行对应的处理




### Links

- Go语言项目开发实战-极客时间


