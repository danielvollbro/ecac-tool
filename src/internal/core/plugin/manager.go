package plugin

import (
	"context"
	"fmt"
	"io"
	"net/rpc"
	"os/exec"

	ecacpluginsdk "github.com/danielvollbro/ecac-plugin-sdk"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
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
	req := ecacpluginsdk.RPCRequest{Params: params}
	var resp ecacpluginsdk.RPCResponse
	err := p.rpc.Call("Plugin.Run", req, &resp)
	if err != nil {
		return "", err
	}
	if resp.Error != "" {
		return "", fmt.Errorf("%s", resp.Error)
	}
	return resp.Result, err
}

func (p *PluginClient) Stop() error {
	err := p.rpc.Close()
	if err != nil {
		return err
	}
	err = p.cmd.Process.Kill()
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

func DecodeTaskConfig(body hcl.Body, plugin HCLPlugin) error {
	ctx := &hcl.EvalContext{}
	diags := gohcl.DecodeBody(body, ctx, plugin.Schema())
	if diags.HasErrors() {
		return fmt.Errorf("HCL decode error: %s", diags.Error())
	}
	return nil
}

func ExecuteTask(ctx context.Context, plugin HCLPlugin) (string, error) {
	if err := plugin.Validate(ctx); err != nil {
		return "", fmt.Errorf("validator failed: %s", err)
	}
	return plugin.Run(ctx)
}
