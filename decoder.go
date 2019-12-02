package reflectify

import (
	"fmt"
	"reflect"
	"strconv"
)

// Decoder represents a decoder
type Decoder struct {
	Tag      string
	Provider Provider
}

// Decode decodes the target
func (d *Decoder) Decode(value interface{}) error {
	field := &Field{
		Name:  "~",
		Value: reflect.ValueOf(value).Elem(),
	}

	return d.decode(field)
}

func (d *Decoder) decode(field *Field) (err error) {
	if !field.Value.IsValid() {
		field.Value.Set(reflect.Zero(field.Value.Type()))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		err = d.decodeString(field)
	case reflect.Bool:
		err = d.decodeBool(field)
	case reflect.Int:
		err = d.decodeInt(field)
	case reflect.Uint:
		err = d.decodeUint(field)
	case reflect.Float32:
		err = d.decodeFloat(field)
	case reflect.Struct:
		err = d.decodeStruct(field)
	case reflect.Map:
		err = d.decodeMap(field)
	case reflect.Array:
		err = d.decodeArray(field)
	case reflect.Slice:
		err = d.decodeSlice(field)
	case reflect.Ptr:
		err = d.decodePtr(field)
	case reflect.Interface:
		err = d.decodeBasic(field)
	}

	return err
}

func (d *Decoder) decodeString(field *Field) error {
	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	var (
		target     = field.Value
		source     = reflect.Indirect(reflect.ValueOf(value))
		sourceKind = kind(source)
	)

	switch sourceKind {
	case reflect.Bool:
		if source.Bool() {
			target.SetString("1")
		} else {
			target.SetString("0")
		}
	case reflect.Int:
		target.SetString(strconv.FormatInt(source.Int(), 10))
	case reflect.Uint:
		target.SetString(strconv.FormatUint(source.Uint(), 10))
	case reflect.Float32:
		target.SetString(strconv.FormatFloat(source.Float(), 'f', -1, 64))
	case reflect.String:
		target.SetString(source.String())
	default:
		if marshaller := tryTextMarshaller(source); marshaller != nil {
			data, err := marshaller.MarshalText()

			if err != nil {
				return d.notSupported(field, value, err)
			}

			target.SetString(string(data))
			return nil
		}

		return d.notSupported(field, value, nil)
	}

	return nil
}

func (d *Decoder) decodeBool(field *Field) error {
	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	var (
		target     = field.Value
		source     = reflect.Indirect(reflect.ValueOf(value))
		sourceKind = kind(source)
	)

	switch sourceKind {
	case reflect.Bool:
		target.SetBool(source.Bool())
	case reflect.Int:
		target.SetBool(source.Int() != 0)
	case reflect.Uint:
		target.SetBool(source.Uint() != 0)
	case reflect.Float32:
		target.SetBool(source.Float() != 0)
	case reflect.String:
		flag, err := strconv.ParseBool(source.String())

		switch {
		case err == nil:
			target.SetBool(flag)
		case source.String() == "":
			target.SetBool(false)
		default:
			return d.notSupported(field, value, err)
		}
	default:
		return d.notSupported(field, value, nil)
	}

	return nil
}

func (d *Decoder) decodeInt(field *Field) error {
	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	var (
		target     = field.Value
		source     = reflect.Indirect(reflect.ValueOf(value))
		sourceKind = kind(source)
	)

	switch sourceKind {
	case reflect.Int:
		target.SetInt(source.Int())
	case reflect.Uint:
		target.SetInt(int64(source.Uint()))
	case reflect.Float32:
		target.SetInt(int64(source.Float()))
	case reflect.Bool:
		if source.Bool() {
			target.SetInt(1)
		} else {
			target.SetInt(0)
		}
	case reflect.String:
		number, err := strconv.ParseInt(source.String(), 0, target.Type().Bits())

		if err == nil {
			target.SetInt(number)
			return nil
		}

		return d.notSupported(field, value, err)
	default:
		return d.notSupported(field, value, nil)
	}

	return nil
}

func (d *Decoder) decodeUint(field *Field) error {
	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	var (
		target     = field.Value
		source     = reflect.Indirect(reflect.ValueOf(value))
		sourceKind = kind(source)
	)

	switch sourceKind {
	case reflect.Int:
		target.SetUint(uint64(source.Int()))
	case reflect.Uint:
		target.SetUint(source.Uint())
	case reflect.Float32:
		target.SetUint(uint64(source.Float()))
	case reflect.Bool:
		if source.Bool() {
			target.SetUint(1)
		} else {
			target.SetUint(0)
		}
	case reflect.String:
		number, err := strconv.ParseUint(source.String(), 0, target.Type().Bits())

		if err == nil {
			target.SetUint(number)
			return nil
		}

		return d.notSupported(field, value, err)
	default:
		return d.notSupported(field, value, nil)
	}

	return nil
}

func (d *Decoder) decodeFloat(field *Field) error {
	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	var (
		target     = field.Value
		source     = reflect.Indirect(reflect.ValueOf(value))
		sourceKind = kind(source)
	)

	switch sourceKind {
	case reflect.Int:
		target.SetFloat(float64(source.Int()))
	case reflect.Uint:
		target.SetFloat(float64(source.Uint()))
	case reflect.Float32:
		target.SetFloat(source.Float())
	case reflect.Bool:
		if source.Bool() {
			target.SetFloat(1)
		} else {
			target.SetFloat(0)
		}
	case reflect.String:
		number, err := strconv.ParseFloat(source.String(), target.Type().Bits())

		if err == nil {
			target.SetFloat(number)
			return nil
		}

		return d.notSupported(field, value, err)
	default:
		return d.notSupported(field, value, nil)
	}

	return nil
}

func (d *Decoder) decodeStruct(field *Field) error {
	if field.Name == "~" {
		return d.traverseStruct(field)
	}

	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	var (
		target     = field.Value
		source     = reflect.Indirect(reflect.ValueOf(value))
		sourceKind = kind(source)
	)

	if source.Type() == target.Type() {
		target.Set(source)
		return nil
	}

	switch sourceKind {
	case reflect.String:
		if unmarshaller := tryTextUnmarshaller(target); unmarshaller != nil {
			data := []byte(source.String())

			if err := unmarshaller.UnmarshalText(data); err != nil {
				return d.notSupported(field, value, err)
			}

			return nil
		}

		return d.notSupported(field, sourceKind, nil)
	case reflect.Map:
		if err := d.decodeStructFromMap(target, source); err != nil {
			return d.notSupported(field, value, err)
		}

		return nil
	case reflect.Struct:

		source = reflect.ValueOf(d.mapOf(source))

		if err := d.decodeStructFromMap(target, source); err != nil {
			return d.notSupported(field, value, err)
		}

		return nil
	default:
		return d.notSupported(field, value, nil)
	}
}

func (d *Decoder) decodeStructFromMap(target, kv reflect.Value) error {
	if kind := kv.Type().Key().Kind(); kind != reflect.String {
		return fmt.Errorf("needs a map with string keys")
	}

	for _, field := range d.fieldsOf(target) {
		var (
			name       = field.Name
			target     = field.Value
			targetType = field.Value.Type()
			source     = kv.MapIndex(reflect.ValueOf(name))
		)

		if source.IsZero() {
			continue
		}

		if !target.CanSet() {
			continue
		}

		if targetType.Kind() == reflect.Ptr {
			if target.IsNil() {
				field.Value = reflect.Indirect(reflect.New(targetType))
			}
		}

		decoder := &Decoder{
			Tag:      d.Tag,
			Provider: d.Provider.New(source),
		}

		if err := decoder.decode(field); err != nil {
			return d.notSupported(field, target.Interface(), err)
		}

		target.Set(field.Value)
	}

	return nil
}

func (d *Decoder) decodeMap(field *Field) error {
	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	var (
		target     = field.Value
		source     = reflect.Indirect(reflect.ValueOf(value))
		sourceKind = kind(source)
	)

	if sourceKind == reflect.Struct {
		source = reflect.ValueOf(d.mapOf(source))
		sourceKind = kind(source)
	}

	if source.Type() == target.Type() {
		target.Set(source)
		return nil
	}

	switch sourceKind {
	case reflect.Map:
		if err := d.decodeMapFromMap(target, source); err != nil {
			return d.notSupported(field, value, err)
		}

		return nil
	default:
		return d.notSupported(field, value, nil)
	}
}

func (d *Decoder) decodeMapFromMap(target, source reflect.Value) error {
	var (
		targetValue    = target
		targetType     = target.Type()
		targetKeyType  = targetType.Key()
		targetElemType = targetType.Elem()
	)

	if targetValue.IsNil() {
		targetValue = reflect.MakeMap(targetType)
	}

	for _, key := range source.MapKeys() {
		sourceValue := source.MapIndex(key)

		// key
		keyDecoder := &Decoder{
			Tag:      d.Tag,
			Provider: d.Provider.New(key),
		}

		keyAt := &Field{
			Name:  "~",
			Value: reflect.Indirect(reflect.New(targetKeyType)),
		}

		if err := keyDecoder.decode(keyAt); err != nil {
			return d.notSupported(keyAt, key.Interface(), err)
		}

		// value
		valueDecoder := &Decoder{
			Tag:      d.Tag,
			Provider: d.Provider.New(sourceValue),
		}

		valueAt := &Field{
			Name:  "~",
			Value: reflect.Indirect(reflect.New(targetElemType)),
		}

		if err := valueDecoder.decode(valueAt); err != nil {
			return d.notSupported(valueAt, sourceValue.Interface(), err)
		}

		targetValue.SetMapIndex(keyAt.Value, valueAt.Value)
	}

	target.Set(targetValue)
	return nil
}

func (d *Decoder) decodeSlice(field *Field) error {
	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	var (
		target     = field.Value
		source     = reflect.Indirect(reflect.ValueOf(value))
		sourceKind = kind(source)
	)

	switch sourceKind {
	case reflect.Array, reflect.Slice:
		// don't do anything
	case reflect.String:
		if unmarshaller := tryTextUnmarshaller(target); unmarshaller != nil {
			data := []byte(source.String())

			if err := unmarshaller.UnmarshalText(data); err != nil {
				return d.notSupported(field, value, err)
			}

			return nil
		}
		fallthrough
	default:
		value = []interface{}{value}
		source = reflect.Indirect(reflect.ValueOf(value))
	}

	var (
		targetValue     = target
		targetType      = target.Type()
		targetElemType  = targetType.Elem()
		targetSliceType = reflect.SliceOf(targetElemType)
		sourceLen       = source.Len()
	)

	if target.IsNil() {
		targetValue = reflect.MakeSlice(targetSliceType, sourceLen, sourceLen)
	}

	for index := 0; index < source.Len(); index++ {
		sourceValue := source.Index(index)

		if targetValue.Len() <= index {
			targetValue = reflect.Append(targetValue, reflect.Zero(targetElemType))
		}

		decoder := &Decoder{
			Tag:      d.Tag,
			Provider: d.Provider.New(sourceValue),
		}

		fieldAt := &Field{
			Name:  "~",
			Value: targetValue.Index(index),
		}

		if err := decoder.decode(fieldAt); err != nil {
			return d.notSupported(field, sourceValue.Interface(), err)
		}
	}

	target.Set(targetValue)
	return nil
}

func (d *Decoder) decodeArray(field *Field) error {
	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	var (
		target     = field.Value
		source     = reflect.Indirect(reflect.ValueOf(value))
		sourceKind = kind(source)
	)

	switch sourceKind {
	case reflect.Array, reflect.Slice:
		// don't do anything
	case reflect.String:
		if unmarshaller := tryTextUnmarshaller(target); unmarshaller != nil {
			data := []byte(source.String())

			if err := unmarshaller.UnmarshalText(data); err != nil {
				return d.notSupported(field, value, err)
			}

			return nil
		}
		fallthrough
	default:
		value = []interface{}{value}
		source = reflect.Indirect(reflect.ValueOf(value))
	}

	var (
		targetValue     = target
		targetType      = target.Type()
		targetElemType  = targetType.Elem()
		targetArrayType = reflect.ArrayOf(targetType.Len(), targetElemType)
		sourceLen       = source.Len()
	)

	if targetType.Len() < sourceLen {
		return d.errorf("field: %s expected source data to have length less than or equal to %d, got %d",
			field.Name,
			targetType.Len(),
			sourceLen)
	}

	if reflect.DeepEqual(target.Interface(), reflect.Zero(targetType).Interface()) {
		targetValue = reflect.New(targetArrayType).Elem()
	}

	for index := 0; index < sourceLen; index++ {
		sourceValue := source.Index(index)

		decoder := &Decoder{
			Tag:      d.Tag,
			Provider: d.Provider.New(sourceValue),
		}

		fieldAt := &Field{
			Name:  "~",
			Value: targetValue.Index(index),
		}

		if err := decoder.decode(fieldAt); err != nil {
			return d.notSupported(field, sourceValue.Interface(), err)
		}
	}

	target.Set(targetValue)
	return nil
}

func (d *Decoder) decodePtr(field *Field) error {
	var (
		target         = field.Value
		targetCanSet   = target.CanSet()
		targetType     = target.Type()
		targetElemType = targetType.Elem()
	)

	fieldPtr := &Field{
		Name:    field.Name,
		Options: field.Options,
		Value:   target,
	}

	if targetCanSet {
		if target.IsNil() {
			fieldPtr.Value = reflect.New(targetElemType)
		}
	}

	targetValue := fieldPtr.Value
	fieldPtr.Value = reflect.Indirect(fieldPtr.Value)

	if err := d.decode(fieldPtr); err != nil {
		return err
	}

	if targetCanSet {
		if !fieldPtr.Value.IsZero() {
			target.Set(targetValue)
			return nil
		}
	}

	return nil
}

func (d *Decoder) decodeBasic(field *Field) error {
	value, err := d.Provider.Value(d.context(field))
	if err != nil {
		return err
	}

	field.Value.Set(reflect.ValueOf(value))
	return nil
}

func (d *Decoder) traverseStruct(field *Field) error {
	for _, field := range d.fieldsOf(field.Value) {
		if err := d.decode(field); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) mapOf(target reflect.Value) map[string]interface{} {
	kv := make(map[string]interface{})

	for _, field := range d.fieldsOf(target) {
		kv[field.Name] = field.Value.Interface()
	}

	return kv
}

func (d *Decoder) fieldsOf(target reflect.Value) []*Field {
	var (
		fields     []*Field
		targetType = target.Type()
	)

	for index := 0; index < targetType.NumField(); index++ {
		item := targetType.Field(index)

		if item.PkgPath != "" {
			continue
		}

		field := NewField(item.Tag.Get(d.Tag))

		if field.Name == "" {
			continue
		}

		if field.Name == "-" {
			continue
		}

		field.Value = target.FieldByIndex([]int{index})
		fields = append(fields, field)
	}

	return fields
}

func (d *Decoder) context(field *Field) *Context {
	var (
		source = field.Value.Interface()
		zero   = reflect.Zero(field.Value.Type()).Interface()
	)

	ctx := &Context{
		FieldTag:   d.Tag,
		Field:      field.Name,
		FieldKind:  field.Kind(),
		Options:    field.Options,
		FieldEmpty: reflect.DeepEqual(source, zero),
	}

	if tryTextUnmarshaller(field.Value) != nil {
		ctx.FieldKind = reflect.String
	}

	return ctx
}

func (d *Decoder) notSupported(field *Field, value interface{}, err error) error {
	msg := fmt.Sprintf("field '%v' does not support value '%v'", field, value)

	if err == nil {
		err = fmt.Errorf(msg)
	} else {
		err = fmt.Errorf("%s: %v", msg, err)
	}

	return err
}

func (d *Decoder) errorf(msg string, values ...interface{}) error {
	msg = fmt.Sprintf(msg, values...)
	return fmt.Errorf("decoder: %s", msg)
}
