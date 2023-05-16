package errors

import (
	"fmt"
	"testing"
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
	MustRegister(coder)
}

func TestWithCode(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	//注册code
	register(ErrUserNotFound, 400, "user not found")
	register(ErrUserEncode, 400, "user  encode error")

	err := WithCode(ErrUserNotFound, "user %s not found.", "cr-mao")
	if err != nil {
		fmt.Println(err.Error()) // user not found
	}
}

func TestWithCodeDetail(t *testing.T) {
	//注册code
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
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
		if IsCode(err, ErrUserEncode) {
			fmt.Println("this is a ErrEncodingFailed error")
		}

		if IsCode(err, ErrUserNotFound) {
			fmt.Println("this is a ErrDatabase error")
		}

		// we can also find the cause error
		fmt.Println(Cause(err))
	}
}

func bindUser() error {
	if err := getUser(); err != nil {
		// Step3: Wrap the error with a new error message and a new error code if needed.
		return WrapC(err, ErrUserEncode, "encoding user %s failed.", "cr-mao")
	}

	return nil
}

func getUser() error {
	if err := queryDatabase(); err != nil {
		// Step2: Wrap the error with a new error message.
		return Wrap(err, "get user failed.")
	}
	return nil
}

func queryDatabase() error {
	// Step1. Create error with specified error code.
	return WithCode(ErrUserNotFound, "user %s not found.", "cr-mao")
}
