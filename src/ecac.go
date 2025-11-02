package main

import (
	"ecac/internal/application/cli"
	"ecac/internal/core/hcl"
	"ecac/internal/infrastructure/storage"
	"encoding/json"
	"fmt"
	"os"
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
	parser := hcl.NewService(*storage)
	config, err := parser.Load(path)
	if err != nil {
		cli.PrintError(err.Error())
		return
	}

	PrintJSON(config)
}

func PrintJSON(obj any) {
	bytes, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(bytes))
}
