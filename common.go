package inflate

import (
	"bytes"
	"database/sql"
	"encoding"
	"fmt"
	"reflect"
	"strings"
)

// Struct works with structs
type Struct struct {
	TagName string
	Value   reflect.Value
}

// StructOf returns the struct
func StructOf(tagName string, value reflect.Value) *Struct {
	return &Struct{
		TagName: tagName,
		Value:   value,
	}
}

// Fields returns the struct fields
func (s *Struct) Fields() []*Field {
	fields := []*Field{}

	for index := 0; index < s.Value.Type().NumField(); index++ {
		var (
			field = s.Value.Type().Field(index)
			value = s.Value.Field(index)
		)

		if field.PkgPath != "" {
			continue
		}

		tag := ParseTag(s.TagName, field.Tag.Get(s.TagName))

		if tag == nil || tag.Name == "-" {
			continue
		}

		if tag.Key != "default" {
			if tag.Name == "" {
				tag.Name = field.Name
			}
		}

		fields = append(fields, &Field{
			Tag:   tag,
			Name:  field.Name,
			Value: value,
		})
	}

	return fields
}

// Map return the struct as map
func (s *Struct) Map() *Map {
	items := make(map[string]interface{})
	s.tree(s, items)

	return MapOf(
		s.TagName,
		reflect.ValueOf(items),
	)
}

func (s *Struct) tree(ch *Struct, kv map[string]interface{}) {
	for _, field := range ch.Fields() {
		if field.Tag.Name == "~" {
			value := elem(field.Value)

			if kind(value) == reflect.Struct {
				s.tree(StructOf(s.TagName, value), kv)
			}

			continue
		}

		if field.Tag.HasOption("omitempty") {
			if field.IsZero() {
				continue
			}
		}

		kv[field.Tag.Name] = field.Value.Interface()
	}
}

// Array return the struct's fields as array
func (s *Struct) Array() *Array {
	var (
		fields = s.Fields()
		items  = make([]interface{}, len(fields))
	)

	for _, field := range fields {
		items = append(items, field.Value.Interface())
	}

	return ArrayOf(
		s.TagName,
		reflect.ValueOf(items),
	)
}

// Tag defines a single struct's string literal tag
type Tag struct {
	Key     string
	Name    string
	Options []string
}

// HasOption returns true if the option is available
func (tag *Tag) HasOption(opt string) bool {
	for _, key := range tag.Options {
		if strings.EqualFold(opt, key) {
			return true
		}
	}

	return false
}

// AddOption adds an option
func (tag *Tag) AddOption(opt string) {
	tag.Options = append(tag.Options, opt)
}

// ParseTag returns the options
func ParseTag(key, value string) *Tag {
	if key == "default" {
		return &Tag{
			Key:  key,
			Name: value,
		}
	}

	parts := strings.Split(value, ",")

	var (
		name string
		opts []string
	)

	if len(parts) == 0 {
		return nil
	}

	if len(parts) > 0 {
		name = parts[0]
	}

	if len(parts) > 1 {
		opts = parts[1:]
	}

	return &Tag{
		Key:     key,
		Name:    name,
		Options: opts,
	}
}

// Field represents a struct field
type Field struct {
	Tag   *Tag
	Name  string
	Value reflect.Value
}

// IsZero return true if it's zero
func (f *Field) IsZero() bool {
	var (
		zero    = reflect.Zero(f.Value.Type()).Interface()
		current = f.Value.Interface()
	)

	return reflect.DeepEqual(current, zero)
}

// Struct returns the field if it's struct
func (f *Field) Struct() *Struct {
	value := elem(f.Value)

	if kind(value) != reflect.Struct {
		return nil
	}

	return &Struct{
		TagName: f.Tag.Key,
		Value:   value,
	}
}

// Map return the struct as map
func (f *Field) Map() *Map {
	if kind(f.Value) != reflect.Map {
		return nil
	}

	return &Map{
		TagName: f.Tag.Key,
		Value:   f.Value,
	}
}

// Array returns the struct as map
func (f *Field) Array() *Array {
	switch kind(f.Value) {
	case reflect.Array, reflect.Slice:
		return &Array{
			TagName: f.Tag.Key,
			Value:   f.Value,
		}
	default:
		return nil
	}
}

// Array works with slices and arrays
type Array struct {
	TagName string
	Elem    reflect.Type
	Value   reflect.Value
}

// MakeArrayOf returns an array / slice
func MakeArrayOf(tagName string, value reflect.Value) *Array {
	array := reflect.MakeSlice(reflect.SliceOf(value.Type()), 1, 1)
	array.Index(0).Set(value)

	return ArrayOf(tagName, array)
}

// ArrayOf returns an array / slice
func ArrayOf(tagName string, value reflect.Value) *Array {
	return &Array{
		TagName: tagName,
		Elem:    value.Type().Elem(),
		Value:   value,
	}
}

// Append appends an item to the array
func (arr *Array) Append(value reflect.Value) {
	target := reflect.New(arr.Elem).Elem()
	set(target, value)

	expanded := reflect.Append(arr.Value, target)
	arr.Value.Set(expanded)
}

// Map represents a map
type Map struct {
	TagName string
	Key     reflect.Type
	Elem    reflect.Type
	Value   reflect.Value
}

// MapOf returns a map
func MapOf(tagName string, value reflect.Value) *Map {
	return &Map{
		TagName: tagName,
		Key:     value.Type().Key(),
		Elem:    value.Type().Elem(),
		Value:   value,
	}
}

// Get returns the value for given key
func (m *Map) Get(key reflect.Value) reflect.Value {
	return m.Value.MapIndex(key)
}

// Values returns the values as array
func (m *Map) Values() *Array {
	var (
		array = reflect.MakeSlice(reflect.SliceOf(m.Elem), 0, 0)
		iter  = m.Value.MapRange()
	)

	for iter.Next() {
		array = reflect.Append(array, iter.Value())
	}

	return &Array{
		TagName: m.TagName,
		Elem:    m.Elem,
		Value:   array,
	}
}

func kind(v reflect.Value) reflect.Kind {
	kind := v.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}

func elem(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		return v.Elem()
	case reflect.Interface:
		return v.Elem()
	}

	return v
}

func refer(v reflect.Value) reflect.Value {
	if !v.IsZero() {
		return elem(v)
	}

	t := v.Type()
	return create(t)
}

func create(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return elem(reflect.New(t))
}

func set(target reflect.Value, source reflect.Value) error {
	for _, value := range variants(source) {
		if value.Type().AssignableTo(target.Type()) {
			target.Set(value)
			return nil
		}
	}

	return rerror(source, target, nil)
}

func variants(source reflect.Value) []reflect.Value {
	values := []reflect.Value{}

	switch source.Type().Kind() {
	case reflect.Ptr:
		values = append(values, source)
		values = append(values, source.Elem())
	default:
		values = append(values, source)

		if source.CanAddr() {
			values = append(values, source.Addr())
		}
	}

	return values
}

func check(name string, value interface{}) (reflect.Value, error) {
	field, ok := value.(reflect.Value)

	if ok {
		return field, nil
	}

	field = reflect.ValueOf(value)

	if field.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("the %v must be a pointer", name)
	}

	if elem := field.Elem(); !elem.CanAddr() {
		return reflect.Value{}, fmt.Errorf("the %v must be addressable (a pointer)", name)
	}

	return elem(field), nil
}

func rerror(source, target reflect.Value, err error) error {
	buffer := &bytes.Buffer{}

	fmt.Fprintf(buffer, "cannot convert %v '%+v' to %v",
		kind(source),
		source.Interface(),
		kind(target),
	)

	if err != nil {
		return fmt.Errorf("%s: %w", buffer.String(), err)
	}

	return fmt.Errorf(buffer.String())
}

func rerrorf(name string, msg interface{}) error {
	return fmt.Errorf("%v: %v", name, msg)
}

func convertable(target reflect.Type) bool {
	var (
		targetInterfaceTypes = []reflect.Type{
			reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem(),
			reflect.TypeOf(new(sql.Scanner)).Elem(),
		}
		targetTypes = []reflect.Type{
			target,
		}
	)

	if target.Kind() != reflect.Ptr {
		targetTypes = append(targetTypes, reflect.PtrTo(target))
	}

	for _, targetInterfaceType := range targetInterfaceTypes {
		for _, targetType := range targetTypes {
			if targetType.Implements(targetInterfaceType) {
				return true
			}
		}
	}

	return false
}

func convertMap(parts []string) (map[string]interface{}, error) {
	var (
		count  = len(parts)
		result = make(map[string]interface{})
	)

	if count%2 != 0 {
		return nil, fmt.Errorf("object value: %s invalid", parts)
	}

	for index := 1; index < count; index = index + 2 {
		prev := index - 1

		var (
			key   = parts[prev]
			value = parts[index]
		)

		if key == "" {
			return nil, fmt.Errorf("object value: %s invalid", parts)
		}

		result[key] = value
	}

	return result, nil
}

func explodeMap(parts []string) (map[string]interface{}, error) {
	var (
		count  = len(parts)
		result = make(map[string]interface{})
	)

	for index := 0; index < count; index++ {
		kv := strings.SplitN(parts[index], "=", 2)

		var (
			key   string
			value string
		)

		switch {
		case len(kv) > 1:
			key = kv[0]
			value = kv[1]
		case len(kv) > 0:
			key = kv[0]
		default:
			return nil, fmt.Errorf("object value: %s invalid", parts[index])
		}

		if key == "" {
			return nil, fmt.Errorf("object value: %s invalid", parts)
		}

		result[key] = value
	}

	return result, nil
}

func convertValue(values []interface{}) interface{} {
	if len(values) == 1 {
		return values[0]
	}

	return values
}

func convertArray(array []string) []interface{} {
	result := make([]interface{}, len(array))

	for index, item := range array {
		result[index] = item
	}

	return result
}
