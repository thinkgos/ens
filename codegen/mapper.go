package codegen

import (
	"fmt"
	"strings"

	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

func (g *CodeGen) GenMapper() *CodeGen {
	if !g.disableDocComment {
		g.Printf("// Code generated by %s. DO NOT EDIT.\n", g.byName)
		g.Printf("// version: %s\n", g.version)
		g.Println()
	}
	g.Println(`syntax = "proto3";`)
	g.Println()
	g.Printf("package %s;\n", g.packageName)
	g.Println()
	if len(g.options) > 0 {
		for k, v := range g.options {
			g.Printf("option %s = \"%s\";\n", k, v)
		}
		g.Println()
	}

	g.Println(`import "protoc-gen-openapiv2/options/annotations.proto";`)
	g.Println()

	for _, et := range g.entities {
		structName := utils.CamelCase(et.Name)

		g.Printf("// %s %s\n", structName, trimStructComment(et.Comment, "\n", "\n// "))

		g.Printf("// #[seaql]\n")
		if et.Table != nil && et.Table.PrimaryKey() != nil {
			g.Printf("// #[seaql(index=\"%s\")]\n", et.Table.PrimaryKey().Definition())
		}
		for _, index := range et.Indexes {
			if et.Table != nil &&
				et.Table.PrimaryKey() != nil &&
				et.Table.PrimaryKey().Index().Name == index.Name {
				continue
			}
			g.Printf("// #[seaql(index=\"%s\")]\n", index.Index.Definition())
		}
		for _, fk := range et.ForeignKeys {
			g.Printf("// #[seaql(foreign_key=\"%s\")]\n", fk.ForeignKey.Definition())
		}
		g.Printf("message %s {\n", structName)
		for i, m := range et.ProtoMessage {
			if m.Comment != "" {
				g.Printf("%s\n", m.Comment)
			}
			g.Println(genMapperMessageField(i+1, m))
		}
		g.Println("}")
	}
	return g
}

func genMapperMessageField(seq int, m *ens.ProtoMessage) string {
	annotation := ""
	if len(m.Annotations) > 0 {
		annotation = fmt.Sprintf(" [%s]", strings.Join(m.Annotations, ", "))
	}
	return fmt.Sprintf("%s %s = %d%s;", m.DataType, m.Name, seq, annotation)
}
