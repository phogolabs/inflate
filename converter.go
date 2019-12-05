package inflate

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

// Converter represents a decoder
type Converter struct {
	TagName string
}

// Convert converts a value to another value
func (d *Converter) Convert(from, to interface{}) error {
	source, err := check("source", from)
	if err != nil {
		return err
	}

	target, err := check("target", to)
	if err != nil {
		return err
	}

	return d.convert(source, target)
}

func (d *Converter) convert(source, target reflect.Value) (err error) {
	if !source.IsValid() {
		source = create(target.Type())
		set(target, source)
		return nil
	}

	if source.Type() == target.Type() {
		set(target, source)
		return nil
	}

	switch target.Kind() {
	case reflect.String:
		err = d.convertToString(source, target)
	case reflect.Bool:
		err = d.convertToBool(source, target)
	case reflect.Int:
		err = d.convertToInt(source, target)
	case reflect.Uint:
		err = d.convertToUint(source, target)
	case reflect.Float32:
		err = d.convertToFloat(source, target)
	case reflect.Struct:
		err = d.convertToStruct(source, target)
	case reflect.Map:
		err = d.convertToMap(source, target)
	case reflect.Array:
		err = d.convertToArray(source, target)
	case reflect.Slice:
		err = d.convertToArray(source, target)
	case reflect.Ptr:
		err = d.convertToPtr(source, target)
	case reflect.Interface:
		err = d.convertToBasic(source, target)
	default:
		return d.error(source, target, nil)
	}

	return err
}

func (d *Converter) convertToString(source, target reflect.Value) error {
	switch source.Kind() {
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
		data, ok, err := d.marshalText(source)
		if ok && err == nil {
			target.SetString(string(data))
			return nil
		}

		return d.error(source, target, err)
	}

	return nil
}

func (d *Converter) convertToBool(source, target reflect.Value) error {
	switch source.Kind() {
	case reflect.Bool:
		target.SetBool(source.Bool())
	case reflect.Int:
		target.SetBool(source.Int() != 0)
	case reflect.Uint:
		target.SetBool(source.Uint() != 0)
	case reflect.Float32:
		target.SetBool(source.Float() != 0)
	case reflect.String:
		value, err := strconv.ParseBool(source.String())

		switch {
		case err == nil:
			target.SetBool(value)
		case source.String() == "":
			target.SetBool(false)
		default:
			return d.error(source, target, err)
		}
	default:
		return d.error(source, target, nil)
	}

	return nil
}

func (d *Converter) convertToInt(source, target reflect.Value) error {
	switch source.Kind() {
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
		value, err := strconv.ParseInt(source.String(), 0, target.Type().Bits())

		if err != nil {
			return d.error(source, target, err)
		}

		target.SetInt(value)
	default:
		return d.error(source, target, nil)
	}

	return nil
}

func (d *Converter) convertToUint(source, target reflect.Value) error {
	switch source.Kind() {
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
		value, err := strconv.ParseUint(source.String(), 0, target.Type().Bits())

		if err != nil {
			return d.error(source, target, err)
		}

		target.SetUint(value)
	default:
		return d.error(source, target, nil)
	}

	return nil
}

func (d *Converter) convertToFloat(source, target reflect.Value) error {
	switch source.Kind() {
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
		value, err := strconv.ParseFloat(source.String(), target.Type().Bits())

		if err != nil {
			return d.error(source, target, err)
		}

		target.SetFloat(value)
	default:
		return d.error(source, target, nil)
	}

	return nil
}

func (d *Converter) convertToStruct(source, target reflect.Value) error {
	switch source.Kind() {
	case reflect.String:
		ok, err := d.unmarshalText(target, []byte(source.String()))
		if ok && err == nil {
			return nil
		}

		return d.error(source, target, err)
	case reflect.Struct:
		return d.convertStructFromMap(
			StructOf(d.TagName, source).Map(),
			StructOf(d.TagName, target),
		)
	case reflect.Map:
		return d.convertStructFromMap(
			MapOf(d.TagName, source),
			StructOf(d.TagName, target),
		)
	default:
		return d.error(source, target, nil)
	}
}

func (d *Converter) convertStructFromMap(source *Map, target *Struct) error {
	for _, field := range target.Fields() {
		key := elem(reflect.New(source.Key))

		if err := d.convert(elem(reflect.ValueOf(field.Tag.Name)), key); err != nil {
			return err
		}

		item := source.Get(key)

		if !item.IsValid() {
			continue
		}

		converted := create(field.Value.Type())

		if err := d.convert(elem(item), converted); err != nil {
			return err
		}

		set(field.Value, converted)
	}

	return nil
}

func (d *Converter) convertToMap(source, target reflect.Value) error {
	switch source.Kind() {
	case reflect.Map:
		return d.convertMapFromMap(
			MapOf(d.TagName, source),
			MapOf(d.TagName, target),
		)
	case reflect.Struct:
		return d.convertMapFromMap(
			StructOf(d.TagName, source).Map(),
			MapOf(d.TagName, target),
		)
	case reflect.String:
		ok, err := d.unmarshalText(target, []byte(source.String()))
		if ok && err == nil {
			return nil
		}

		return d.error(source, target, err)
	default:
		return d.error(source, target, nil)
	}
}

func (d *Converter) convertMapFromMap(source *Map, target *Map) error {
	if source.Value.Type() == target.Value.Type() {
		target.Value.Set(source.Value)
		return nil
	}

	iter := source.Value.MapRange()

	for iter.Next() {
		key := elem(reflect.New(target.Key))

		if err := d.convert(elem(iter.Key()), key); err != nil {
			return err
		}

		var (
			item      = iter.Value()
			converted = create(target.Elem)
		)

		if err := d.convert(elem(item), converted); err != nil {
			return err
		}

		if !converted.IsZero() {
			target.Value.SetMapIndex(key, converted)
		}
	}

	return nil
}

func (d *Converter) convertToArray(source, target reflect.Value) error {
	switch source.Kind() {
	case reflect.String:
		if ok, err := d.unmarshalText(target, []byte(source.String())); ok {
			if err != nil {
				return d.error(source, target, err)
			}

			return nil
		}
	case reflect.Array, reflect.Slice:
		return d.convertArrayFromArray(
			ArrayOf(d.TagName, source),
			ArrayOf(d.TagName, target),
		)
	case reflect.Map:
		return d.convertArrayFromArray(
			MapOf(d.TagName, source).Values(),
			ArrayOf(d.TagName, target),
		)
	case reflect.Struct:
		return d.convertArrayFromArray(
			StructOf(d.TagName, source).Array(),
			ArrayOf(d.TagName, target),
		)
	}

	return d.convertArrayFromArray(
		MakeArrayOf(d.TagName, source),
		ArrayOf(d.TagName, target),
	)
}

func (d *Converter) convertArrayFromArray(source *Array, target *Array) error {
	for index := 0; index < source.Value.Len(); index++ {
		var (
			item      = elem(source.Value.Index(index))
			converted = create(target.Elem)
		)

		if err := d.convert(item, converted); err != nil {
			return err
		}

		if !converted.IsZero() {
			switch target.Value.Kind() {
			case reflect.Array:
				if index >= target.Value.Len() {
					return nil
				}

				set(target.Value.Index(index), converted)
			case reflect.Slice:
				target.Append(converted)
			}
		}
	}

	return nil
}

func (d *Converter) convertToPtr(source, target reflect.Value) error {
	origin := target

	if target.IsNil() {
		target = create(target.Type().Elem())
	} else {
		target = elem(target)
	}

	if err := d.convert(elem(source), target); err != nil {
		return err
	}

	set(origin, target)
	return nil
}

func (d *Converter) convertToBasic(source, target reflect.Value) error {
	origin := target

	if !source.IsValid() {
		source = create(source.Type())
	}

	if target.IsValid() && target.Elem().IsValid() {
		target = create(target.Elem().Type())
		source = elem(source)

		if err := d.convert(source, target); err != nil {
			return err
		}

		set(origin, target)
		return nil
	}

	if source.Type().AssignableTo(target.Type()) {
		set(target, source)
		return nil
	}

	if source.CanAddr() {
		if source.Addr().Type().AssignableTo(target.Type()) {
			set(target, source.Addr())
			return nil
		}
	}

	return d.error(source, target, nil)
}

func (d *Converter) marshalText(source reflect.Value) ([]byte, bool, error) {
	targetType := reflect.TypeOf(new(encoding.TextMarshaler)).Elem()

	if source.Type().Implements(targetType) {
		var (
			encoder   = source.Interface().(encoding.TextMarshaler)
			data, err = encoder.MarshalText()
		)

		return data, true, err
	}

	return nil, false, nil
}

func (d *Converter) unmarshalText(target reflect.Value, data []byte) (bool, error) {
	targetType := reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

	if target.Kind() != reflect.Ptr {
		if target.CanAddr() {
			target = target.Addr()
		}
	}

	if target.Type().Implements(targetType) {
		var (
			decoder = target.Interface().(encoding.TextUnmarshaler)
			err     = decoder.UnmarshalText(data)
		)

		return true, err
	}

	return false, nil
}

func (d *Converter) error(source, target reflect.Value, err error) error {
	buffer := &bytes.Buffer{}

	fmt.Fprintf(buffer, "cannot convert %v '%+v' to %v",
		source.Kind(),
		source.Interface(),
		target.Kind())

	if err != nil {
		return fmt.Errorf("%s: %w", buffer.String(), err)
	}

	return fmt.Errorf(buffer.String())
}
