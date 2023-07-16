// Code generated by internal/integer.tpl, DO NOT EDIT.

package ens

import (
	"reflect"
)

var _ Fielder = (*uintBuilder)(nil)
var uintType = reflect.TypeOf(uint(0))

// Uint returns a new Field with type uint.
func Uint(name string) *uintBuilder {
	return &uintBuilder{
		&FieldDescriptor{
			Name: name,
			Type: UintType(),
		},
	}
}

// uintBuilder is the builder for uint field.
type uintBuilder struct {
	inner *FieldDescriptor
}

// Comment sets the comment of the field.
func (b *uintBuilder) Comment(c string) *uintBuilder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *uintBuilder) Nullable() *uintBuilder {
	b.inner.Nullable = true
	return b
}

// Definition set the sql definition of the field.
func (b *uintBuilder) Definition(s string) *uintBuilder {
	b.inner.Definition = s
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Uint("uint").
//		GoType(pkg.Uint(0))
func (b *uintBuilder) GoType(typ any) *uintBuilder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *uintBuilder) Optional() *uintBuilder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of tags to the field tag.
//
//	field.Uint("uint").
//		Tags("yaml:"xxx"")
func (b *uintBuilder) Tags(tags ...string) *uintBuilder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Fielder interface by returning its descriptor.
func (b *uintBuilder) Build(opt *Option) *FieldDescriptor {
	//	b.inner.checkGoType(uintType)
	b.inner.build(opt)
	return b.inner
}
