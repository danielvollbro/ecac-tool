package main

import (
	"ecac/internal/application/cli"
	hclService "ecac/internal/core/hcl"
	"ecac/internal/core/plugin"
	"ecac/internal/infrastructure/storage"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func main() {
	cli := cli.Cli{
		Level: 2,
	}
	if len(os.Args[1:]) <= 0 {
		cli.PrintError("you need to specify a file")
		return
	}

	path := os.Args[1]

	storage := storage.NewService()
	parser := hclService.NewService(*storage)
	config, err := parser.ParseConfig(path)
	if err != nil {
		cli.PrintError(err.Error())
		return
	}

	activePlugins := map[string]*plugin.PluginClient{}
	for _, p := range config.Plugins {
		fmt.Printf("Fetching plugin %q -> %s@%s\n", p.Name, p.Source, p.Version)
		path, err := plugin.FetchAndBuild(p.Source, p.Version)
		if err != nil {
			panic(err)
		}

		client, err := plugin.StartPlugin(path)
		if err != nil {
			panic(err)
		}
		activePlugins[p.Name] = client
	}

	for _, t := range config.Tasks {
		client := activePlugins[t.Plugin]
		if client == nil {
			fmt.Printf("Plugin %q not loaded, skipping task %q\n", t.Plugin, t.Name)
			continue
		}

		params, err := decodeTaskParams(t.Config)
		if err != nil {
			fmt.Printf("Failed to decode config for task %q: %v\n", t.Name, err)
			continue
		}

		result, err := client.Run(params)
		if err != nil {
			fmt.Printf("Task %q failed: %v\n", t.Name, err)
			continue
		}

		fmt.Printf("Task %q completed: %s\n", t.Name, result)
	}

	for _, c := range activePlugins {
		_ = c.Stop()
	}

	// PrintJSON(config)
}

func PrintJSON(obj any) {
	bytes, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(bytes))
}

func decodeTaskParams(body hcl.Body) (map[string]any, error) {
	type Raw struct {
		Remain hcl.Body `hcl:",remain"`
	}
	var raw Raw
	ctx := &hcl.EvalContext{}
	diags := gohcl.DecodeBody(body, ctx, &raw)
	if diags.HasErrors() {
		return nil, diags
	}

	schema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "*"},
		},
	}
	content, _, _ := raw.Remain.PartialContent(schema)

	out := make(map[string]any)
	for name, attr := range content.Attributes {
		val, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			continue
		}
		out[name] = ctyToGo(val)
	}

	return out, nil
}

func ctyToGo(v cty.Value) any {
	if !v.IsKnown() || v.IsNull() {
		return nil
	}

	if v.Type() == cty.String {
		var s string
		err := gocty.FromCtyValue(v, &s)
		if err != nil {
			panic(err)
		}
		return s
	}

	if v.Type() == cty.Bool {
		var b bool
		err := gocty.FromCtyValue(v, &b)
		if err != nil {
			panic(err)
		}
		return b
	}

	if v.Type() == cty.Number {
		var f float64
		err := gocty.FromCtyValue(v, &f)
		if err != nil {
			panic(err)
		}
		return f
	}

	if v.Type().IsListType() || v.Type().IsTupleType() {
		var list []cty.Value
		err := gocty.FromCtyValue(v, &list)
		if err != nil {
			panic(err)
		}
		res := make([]any, 0, len(list))
		for _, elem := range list {
			res = append(res, ctyToGo(elem))
		}
		return res
	}

	if v.Type().IsMapType() || v.Type().IsObjectType() {
		m := map[string]any{}
		for k, val := range v.AsValueMap() {
			m[k] = ctyToGo(val)
		}
		return m
	}

	return fmt.Sprintf("%v", v.GoString())
}
