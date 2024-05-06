package rapier

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/things-go/ens"
	"golang.org/x/tools/imports"
)

type CodeGen struct {
	buf               bytes.Buffer
	Structs           []*Struct // required
	PackageName       string    // required
	ModelImportPath   string    // required
	ByName            string
	Version           string
	DisableDocComment bool
}

// Bytes returns the CodeBuf's buffer.
func (g *CodeGen) Bytes() []byte {
	return g.buf.Bytes()
}

// FormatSource return formats and adjusts imports contents of the CodeGen's buffer.
func (g *CodeGen) FormatSource() ([]byte, error) {
	data := g.buf.Bytes()
	if len(data) == 0 {
		return data, nil
	}
	// return format.Source(data)
	return imports.Process("", data, nil)
}

// Write appends the contents of p to the buffer,
func (g *CodeGen) Write(b []byte) (n int, err error) {
	return g.buf.Write(b)
}

// Print formats using the default formats for its operands and writes to the generated output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Print(a ...any) (n int, err error) {
	return fmt.Fprint(&g.buf, a...)
}

// Printf formats according to a format specifier for its operands and writes to the generated output.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(&g.buf, format, a...)
}

// Fprintln formats using the default formats to the generated output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Println(a ...any) (n int, err error) {
	return fmt.Fprintln(&g.buf, a...)
}

func (g *CodeGen) Get() *CodeGen {
	pkgQualifierPrefix := ""
	if p := ens.PkgName(g.ModelImportPath); p != "" {
		pkgQualifierPrefix = p + "."
	}
	if !g.DisableDocComment {
		g.Printf("// Code generated by %s. DO NOT EDIT.\n", g.ByName)
		g.Printf("// version: %s\n", g.Version)
		g.Println()
	}
	g.Printf("package %s\n", g.PackageName)
	g.Println()

	//* import
	g.Println("import (")
	if pkgQualifierPrefix != "" {
		g.Printf("\"%s\"\n", g.ModelImportPath)
		g.Println()
	}
	g.Println(`rapier "github.com/thinkgos/gorm-rapier"`)
	g.Println(`"gorm.io/gorm"`)
	g.Println(")")

	//* struct
	for _, et := range g.Structs {
		structGoName := et.GoName
		tableName := et.TableName

		varRefNative := fmt.Sprintf(`ref_%s_Native`, structGoName)
		funcInnerNew := fmt.Sprintf(`new_%s`, structGoName)
		{ //* var field
			g.Printf("var %s = New_%s(\"%s\")\n", varRefNative, structGoName, tableName)
			g.Println()
		}

		typeNative := fmt.Sprintf("%s_Native", structGoName)
		//* type
		{
			g.Printf("type %s struct {\n", typeNative)
			g.Println("refAlias string")
			g.Println("refTableName string")
			g.Println("ALL rapier.Asterisk")
			for _, field := range et.Fields {
				g.Printf("%s rapier.%s\n", field.GoName, field.Type.String())
			}
			g.Println("}")
			g.Println()
		}
		//* function new_xxx
		{
			g.Printf("func %s(tableName, alias string) *%s {\n", funcInnerNew, typeNative)
			g.Printf("return &%s {\n", typeNative)
			g.Println("refAlias: alias,")
			g.Println("refTableName: tableName,")
			g.Println("ALL:  rapier.NewAsterisk(alias),")
			for _, field := range et.Fields {
				g.Printf("%s: rapier.New%s(alias, \"%s\"),\n", field.GoName, field.Type.String(), field.ColumnName)
			}
			g.Println("}")
			g.Println("}")
			g.Println()
		}
		//* function Ref_xxx
		{
			g.Printf("// Ref_%s model with TableName `%s`.\n", structGoName, tableName)
			g.Println("// NOTE: Don't modify any public field!!!")
			g.Printf("func Ref_%s() *%s {\n", structGoName, typeNative)
			g.Printf("return %s\n", varRefNative)
			g.Println("}")
			g.Println()
		}
		//* function New_xxxx
		{
			g.Printf("// New_%s new instance.\n", structGoName)
			g.Println("// NOTE: Don't modify any public field!!!")
			g.Printf("func New_%s(tableName string) *%s {\n", structGoName, typeNative)
			g.Printf("return %s(tableName, tableName)\n", funcInnerNew)
			g.Println("}")
			g.Println()
		}
		//* method As
		{
			g.Println("// As alias")
			g.Printf("func (x *%[1]s) As(alias string) *%[1]s {\n", typeNative)
			g.Printf("return %s(x.refTableName, alias)\n", funcInnerNew)
			g.Println("}")
			g.Println()
		}
		//* method Ref_Alias
		{
			g.Printf("// Ref_Alias hold alias name when call %[1]s_Native.As that you defined.\n", structGoName)
			g.Printf("func (x *%s) Ref_Alias() string {\n", typeNative)
			g.Println("return x.refAlias")
			g.Println("}")
			g.Println()
		}
		// impl TableName interface
		{
			//* method TableName
			g.Printf("// TableName hold table name when call New_%[1]s that you defined.\n", structGoName)
			g.Printf("func (x *%s) TableName() string {\n", typeNative)
			g.Println("return x.refTableName")
			g.Println("}")
			g.Println()
		}

		//* method New_Executor
		{
			modelName := pkgQualifierPrefix + structGoName
			g.Println("// New_Executor new entity executor which suggest use only once.")
			g.Printf("func (*%s) New_Executor(db *gorm.DB) *rapier.Executor[%s] {\n", typeNative, modelName)
			g.Printf("return rapier.NewExecutor[%s](db)\n", modelName)
			g.Println("}")
			g.Println()
		}
		//* method Select_Expr
		{
			g.Println("// Select_Expr select model fields")
			g.Printf("func (x *%s) Select_Expr() []rapier.Expr {\n", typeNative)
			g.Println("return []rapier.Expr{")
			for _, field := range et.Fields {
				g.Printf("x.%s,\n", field.GoName)
			}
			g.Println("}")
			g.Println("}")
			g.Println()
		}

		//* method Select_VariantExpr
		{
			g.Println("// Select_VariantExpr select model fields, but time.Time field convert to timestamp(int64).")
			g.Printf("func (x *%s) Select_VariantExpr(prefixes ...string) []rapier.Expr {\n", typeNative)
			g.Println("if len(prefixes) > 0 && prefixes[0] != \"\" {")
			g.Println("return []rapier.Expr{")
			for _, field := range et.Fields {
				g.Println(genRapier_SelectVariantExprField(field, true))
			}
			g.Println("}")
			g.Println("} else {")
			g.Println("return []rapier.Expr{")
			for _, field := range et.Fields {
				g.Println(genRapier_SelectVariantExprField(field, false))
			}
			g.Println("}")
			g.Println("}")
			g.Println("}")
			g.Println()
		}
	}
	return g
}

func genRapier_SelectVariantExprField(field *StructField, hasPrefix bool) string {
	goName := field.GoName

	b := &strings.Builder{}
	b.Grow(64)
	b.WriteString("x.")
	b.WriteString(goName)
	if field.Type == Time {
		b.WriteString(".UnixTimestamp()")
		if field.Nullable {
			b.WriteString(".IfNull(0)")
		}
		if !hasPrefix {
			fmt.Fprintf(b, ".As(x.%s.ColumnName())", goName)
		}
	}
	if hasPrefix {
		fmt.Fprintf(b, ".As(x.%s.FieldName(prefixes...))", goName)
	}
	b.WriteString(",")
	return b.String()
}
