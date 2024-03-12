package internal

import (
	"fmt"
	natspb "github.com/mrdan4es/bazel-monorepo-example/api/nats/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".pb.nats.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// Code generated by protoc-gen-nats. DO NOT EDIT.")
	g.P("// source: ", file.Desc.Path())
	g.P()
	g.P("package ", file.GoPackageName, "// ", file.GoImportPath)
	g.P()

	generate(file, g)

	return g
}

func generate(f *protogen.File, g *protogen.GeneratedFile) {
	if len(f.Services) == 0 {
		return
	}

	for _, service := range f.Services {
		serviceInfo, ok := proto.GetExtension(service.Desc.Options(), natspb.E_NatsService).(*natspb.NatsService)
		if !ok {
			continue
		}

		for _, method := range service.Methods {
			methodInfo, ok := proto.GetExtension(method.Desc.Options(), natspb.E_NatsMethod).(*natspb.NatsMethod)
			if !ok {
				continue
			}
			g.P(fmt.Sprintf("func (*c NatsClient) %s(res *%s) error {", method.GoName, method.Input.GoIdent.GoName))
			g.P(fmt.Sprintf("// topic pub: %s.%s", serviceInfo.Topic, methodInfo.ResponseName))
			g.P(fmt.Sprintf("// topic sub: %s.%s", serviceInfo.Topic, methodInfo.RequestName))
			g.P("return c.Send(%s, res.Bytes())")
			g.P("}")
		}
	}
}
