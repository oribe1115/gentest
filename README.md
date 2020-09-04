# gentest

`gentest` generate test code template for function

## Installation

```sh
$ go get -u github.com/oribe1115/gentest/cmd/gentest
```

## Example

```go
package main

import "fmt"

func main() {
	greeting := hello("world")
	fmt.Println(greeting)
}

func hello(s string) string {
	return "Hello " + s
}
```

```sh
$ go vet -vettool=`which gentest` -gentest.offset=101 main.go

func TestHello(t *testing.T) {

        type input struct {
                s string
        }
        type expected struct {
                gotString string
        }
        tests := []struct {
                Name     string
                Input    input
                Expected expected
        }{
                // TODO: Add test cases.
        }

        for _, test := range tests {
                t.Run(test.Name, func(t *testing.T) {

                        gotString := hello(test.Input.s)

                        assert.Equal(t, test.Expected.gotString, gotString)
                })
        }
}
```

## Flags

### gentest.offset
Essential flag
Target function is searched by this
Default: `0`

### gentest.paralell
Optional flag
If parallel is `true`, generated test code contains t.Parallel and t.Cleanup
Default: `false`

```go
package main

import "fmt"

func main() {
	greeting := hello("world")
	fmt.Println(greeting)
}

func hello(s string) string {
	return "Hello " + s
}
```

```sh
go vet -vettool=`which gentest` -gentest.offset=101 -gentest.parallel=true main.go

func TestHello(t *testing.T) {
        t.Parallel()
        type input struct {
                s string
        }
        type expected struct {
                gotString string
        }
        tests := []struct {
                Name     string
                Input    input
                Expected expected
        }{
                // TODO: Add test cases.
        }

        for _, test := range tests {
                t.Run(test.Name, func(t *testing.T) {
                        t.Parallel()
                        t.Cleanup()

                        gotString := hello(test.Input.s)

                        assert.Equal(t, test.Expected.gotString, gotString)
                })
        }
}
```

## Pattern
`gentest` generate best test code template for target function

### Retrun values contin error
```go
package main

import (
	"errors"
)

func main() {
	_, _ = errFunc()
}

func errFunc() (string, error) {
	return "", errors.New("error")
}
```

```sh
$ go vet -vettool=`which gentest` -gentest.offset=76 main.go

func TestErrFunc(t *testing.T) {

        type expected struct {
                gotString string
                gotError  error
        }
        tests := []struct {
                Name      string
                Expected  expected
                WantError bool
        }{
                // TODO: Add test cases.
        }

        for _, test := range tests {
                t.Run(test.Name, func(t *testing.T) {

                        gotString, gotError := errFunc()

                        if test.WantError {
                                assert.Error(t, gotError)
                                if test.Expected.gotError != nil {
                                        assert.EqualError(t, gotError, test.Expected.gotError.String())
                                }
                        } else {
                                assert.NoError(t, gotError)
                        }

                        assert.Equal(t, test.Expected.gotString, gotString)
                })
        }
}
```

### Target function is method
```go
package main

type T struct {
	Hoge string
}

func main() {
	t := T{
		Hoge: "hoge",
	}
	t.method()
}

func (t T) method() string {
	return t.Hoge
}
```

```sh
go vet -vettool=`which gentest` -gentest.offset=114 main.go

func TestMethod(t *testing.T) {

        type expected struct {
                gotString string
        }
        tests := []struct {
                Name     string
                Use      main.T
                Expected expected
        }{
                // TODO: Add test cases.
        }

        for _, test := range tests {
                t.Run(test.Name, func(t *testing.T) {

                        gotString := test.Use.method()

                        assert.Equal(t, test.Expected.gotString, gotString)
                })
        }
}
```

### Method change receiver field
```go
package main

type T struct {
	Hoge string
}

func main() {
	t := &T{
		Hoge: "hoge",
	}
	t.method()
}

func (t *T) method() string {
	t.Hoge = "fuga"
	return t.Hoge
}
```

```sh
go vet -vettool=`which gentest` -gentest.offset=116 main.go

func TestMethod(t *testing.T) {

        type expected struct {
                gotString string
        }
        tests := []struct {
                Name        string
                Use         *main.T
                Expected    expected
                UseExpected *main.T
        }{
                // TODO: Add test cases.
        }

        for _, test := range tests {
                t.Run(test.Name, func(t *testing.T) {

                        gotString := test.Use.method()

                        assert.Equal(t, test.Expected.gotString, gotString)
                        assert.Equal(t, test.UseExpected, test.Use)
                })
        }
}
```