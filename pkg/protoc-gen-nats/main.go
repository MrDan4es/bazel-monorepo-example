package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/mrdan4es/bazel-monorepo-example/pkg/protoc-gen-nats/internal"
)

func main() {
	flag.Parse()

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if f.Generate {
				internal.GenerateFile(gen, f)
			}
		}

		return nil
	})
}
