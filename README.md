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

## フラグ

- `offset`: 対象となる関数のオフセット 必須
- `paralell`: trueにセットすると`t.Parallel`を追加したテストコードを生成する　オプション
