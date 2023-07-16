// Code generated by internal/integer.tpl, DO NOT EDIT.

package ens

import (
	"reflect"
)

var _ Fielder = (*uint8Builder)(nil)
var uint8Type = reflect.TypeOf(uint8(0))

// Uint8 returns a new Field with type uint8.
func Uint8(name string) *uint8Builder {
	return &uint8Builder{
		&FieldDescriptor{
			Name: name,
			Type: Uint8Type(),
		},
	}
}

// uint8Builder is the builder for uint8 field.
type uint8Builder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *uint8Builder) Comment(c string) *uint8Builder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *uint8Builder) Nullable() *uint8Builder {
	b.inner.Nullable = true
	return b
}

// Definition set the sql definition of the field.
func (b *uint8Builder) Definition(s string) *uint8Builder {
	b.inner.Definition = s
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Uint8("uint8").
//		GoType(pkg.Uint8(0))
func (b *uint8Builder) GoType(typ any) *uint8Builder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *uint8Builder) Optional() *uint8Builder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of tags to the field tag.
//
//	field.Uint8("uint8").
//		Tags("yaml:"xxx"")
func (b *uint8Builder) Tags(tags ...string) *uint8Builder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Fielder interface by returning its descriptor.
func (b *uint8Builder) Build(opt *Option) *FieldDescriptor {
	//	b.inner.checkGoType(uint8Type)
	b.inner.build(opt)
	return b.inner
}
