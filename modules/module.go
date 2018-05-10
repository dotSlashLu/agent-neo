package module

type Module interface {
	Call(fn string, params []byte) (string, error)
}
