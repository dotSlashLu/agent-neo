package module

// import "github.com/dotSlashLu/agent-neo/lib"

type CallableModule interface {
	Call(fn string, params []byte) (string, error)
}
