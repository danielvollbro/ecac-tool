package plugin

import (
	"io"
	"net/rpc"
	"os/exec"
)

type PluginClient struct {
	cmd *exec.Cmd
	rpc *rpc.Client
}

func StartPlugin(path string) (*PluginClient, error) {
	cmd := exec.Command(path, "serve")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	client := rpc.NewClient(&rw{r: stdout, w: stdin})
	return &PluginClient{cmd: cmd, rpc: client}, nil
}

func (p *PluginClient) Run(params map[string]any) (string, error) {
	var result string
	err := p.rpc.Call("Plugin.Run", params, &result)
	return result, err
}

func (p *PluginClient) Stop() error {
	err := p.rpc.Close()
	if err != nil {
		return err
	}

	return p.cmd.Wait()
}

type rw struct {
	r io.Reader
	w io.Writer
}

func (r *rw) Read(b []byte) (int, error) { return r.r.Read(b) }

func (r *rw) Write(b []byte) (int, error) { return r.w.Write(b) }

func (r *rw) Close() error { return nil }
