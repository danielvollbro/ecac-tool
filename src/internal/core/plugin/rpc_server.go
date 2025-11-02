package plugin

import (
	"encoding/json"
	"net/rpc"
	"os"
)

type PluginServer interface {
	Run(params map[string]any) (string, error)
}

func Server(p PluginServer) {
	rpc.RegisterName("Plugin", p)
	rpc.ServeConn(&rw{r: os.Stdin, w: os.Stdout})
}

func EncodeParams(args []string) map[string]any {
	m := map[string]any{}
	for _, a := range args {
		var kv map[string]any
		json.Unmarshal([]byte(a), &kv)
		for k, v := range kv {
			m[k] = v
		}
	}
	return m
}
