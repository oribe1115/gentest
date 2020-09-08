# gentest

`gentest` generates test code template for function
Generated test code uses [testify](github.com/stretchr/testify)

## Installation

```sh
$ go get -u github.com/oribe1115/gentest/cmd/gentest
```

## Example

```go:sample.go
package sample

func hello(s string) string {
	return "Hello " + s
}
```

```sh
$ go vet -vettool=`which gentest` -gentest.offset=21 sample.go

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
                test := test
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

```go:sample.go
package sample

func hello(s string) string {
	return "Hello " + s
}
```

```sh
go vet -vettool=`which gentest` -gentest.offset=21 -gentest.parallel=true sample.go

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
                test := test
                t.Run(test.Name, func(t *testing.T) {
                        t.Parallel()
                        t.Cleanup(
                        // TODO: Add function.
                        )
                        gotString := hello(test.Input.s)

                        assert.Equal(t, test.Expected.gotString, gotString)
                })
        }
}
```

## Pattern
`gentest` generates best test code template for target function

### Retrun values contin error
```go:sample.go
package sample

import "errors"

func errFunc() (string, error) {
	return "", errors.New("error")
}
```

```sh
$ go vet -vettool=`which gentest` -gentest.offset=40 sample.go

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
                test := test
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
```go:sample.go
package sample

type T struct {
	Hoge string
}

func (t T) Method() string {
	return t.Hoge
}
```

```sh
$ go vet -vettool=`which gentest` -gentest.offset=114 sample.go

func TestT_Method(t *testing.T) {

        type expected struct {
                gotString string
        }
        tests := []struct {
                Name     string
                Use      sample.T
                Expected expected
        }{
                // TODO: Add test cases.
        }

        for _, test := range tests {
                test := test
                t.Run(test.Name, func(t *testing.T) {

                        gotString := test.Use.Method()

                        assert.Equal(t, test.Expected.gotString, gotString)
                })
        }
}
```

### Method change receiver field
```go:sample.go
package sample

type T struct {
	Hoge string
}

func (t *T) Method() string {
	t.Hoge = "hoge"
	return t.Hoge
}
```

```sh
$ go vet -vettool=`which gentest` -gentest.offset=116 sample.go

func TestT_Method(t *testing.T) {

        type expected struct {
                gotString string
        }
        tests := []struct {
                Name        string
                Use         *sample.T
                Expected    expected
                UseExpected *sample.T
        }{
                // TODO: Add test cases.
        }

        for _, test := range tests {
                test := test
                t.Run(test.Name, func(t *testing.T) {

                        gotString := test.Use.Method()

                        assert.Equal(t, test.Expected.gotString, gotString)
                        assert.Equal(t, test.UseExpected, test.Use)
                })
        }
}
```