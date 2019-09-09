//go:generate msgp
package api

const (
	SayHello = "SayHello"
)

type HelloRequest struct {
	Name string
}

type HelloReply struct {
	Message string
}
