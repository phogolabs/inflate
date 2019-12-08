package inflate

import (
	"database/sql"
	"database/sql/driver"
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
		source = refer(target)
		return set(target, source)
	}

	if source.Type() == target.Type() {
		return set(target, source)
	}

	switch kind(target) {
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
		return rerror(source, target, nil)
	}

	return err
}

func (d *Converter) convertToString(source, target reflect.Value) error {
	switch kind(source) {
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
		data, ok, err := d.textMarshal(source)
		if ok && err == nil {
			target.SetString(data)
			return nil
		}

		value, ok, err := d.valueRead(source)
		if ok && err == nil {
			source = elem(reflect.ValueOf(value))
			return d.convertToString(source, target)
		}

		return rerror(source, target, err)
	}

	return nil
}

func (d *Converter) convertToBool(source, target reflect.Value) error {
	switch kind(source) {
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
			return rerror(source, target, err)
		}
	case reflect.Struct:
		value, ok, err := d.valueRead(source)
		if ok && err == nil {
			source = elem(reflect.ValueOf(value))
			return d.convertToBool(source, target)
		}

		return rerror(source, target, err)
	default:
		return rerror(source, target, nil)
	}

	return nil
}

func (d *Converter) convertToInt(source, target reflect.Value) error {
	switch kind(source) {
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
			return rerror(source, target, err)
		}

		target.SetInt(value)
	case reflect.Struct:
		value, ok, err := d.valueRead(source)
		if ok && err == nil {

			source = elem(reflect.ValueOf(value))
			return d.convertToInt(source, target)
		}

		return rerror(source, target, err)
	default:
		return rerror(source, target, nil)
	}

	return nil
}

func (d *Converter) convertToUint(source, target reflect.Value) error {
	switch kind(source) {
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
			return rerror(source, target, err)
		}

		target.SetUint(value)
	case reflect.Struct:
		value, ok, err := d.valueRead(source)
		if ok && err == nil {
			source = elem(reflect.ValueOf(value))
			return d.convertToUint(source, target)
		}

		return rerror(source, target, err)
	default:
		return rerror(source, target, nil)
	}

	return nil
}

func (d *Converter) convertToFloat(source, target reflect.Value) error {
	switch kind(source) {
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
			return rerror(source, target, err)
		}

		target.SetFloat(value)
	case reflect.Struct:
		value, ok, err := d.valueRead(source)
		if ok && err == nil {
			source = elem(reflect.ValueOf(value))
			return d.convertToFloat(source, target)
		}

		return rerror(source, target, err)
	default:
		return rerror(source, target, nil)
	}

	return nil
}

func (d *Converter) convertToStruct(source, target reflect.Value) error {
	switch kind(source) {
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
	case reflect.String:
		ok, err := d.textUnmarshal(source.String(), target)
		if ok && err == nil {
			return nil
		}

		fallthrough
	default:
		ok, err := d.valueScan(source.Interface(), target)
		if ok && err == nil {
			return nil
		}

		return rerror(source, target, err)
	}
}

func (d *Converter) convertStructFromMap(source *Map, target *Struct) error {
	for _, field := range target.Fields() {
		if field.Tag.Name == "~" {
			value := refer(field.Value)

			switch kind(value) {
			case reflect.Struct:
				obj := StructOf(d.TagName, value)

				if err := d.convertStructFromMap(source, obj); err != nil {
					return rerrorf(field.Name, err)
				}
			case reflect.Map:
				obj := MapOf(d.TagName, value)

				if err := d.convertMapFromMap(source, obj); err != nil {
					return rerrorf(field.Name, err)
				}
			default:
				return rerrorf(field.Name, rerror(source.Value, target.Value, nil))
			}

			if err := set(field.Value, value.Addr()); err != nil {
				return err
			}

			continue
		}

		key := elem(reflect.New(source.Key))

		if err := d.convert(elem(reflect.ValueOf(field.Tag.Name)), key); err != nil {
			return err
		}

		item := source.Get(key)

		if !item.IsValid() {
			continue
		}

		converted := refer(field.Value)

		if err := d.convert(elem(item), converted); err != nil {
			return err
		}

		if err := set(field.Value, converted); err != nil {
			return err
		}
	}

	return nil
}

func (d *Converter) convertToMap(source, target reflect.Value) error {
	switch kind(source) {
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
		ok, err := d.textUnmarshal(source.String(), target)
		if ok && err == nil {
			return nil
		}

		fallthrough
	default:
		ok, err := d.valueScan(source.Interface(), target)
		if ok && err == nil {
			return nil
		}

		return rerror(source, target, nil)
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
			return rerrorf(fmt.Sprintf("%v", iter.Key().Interface()), err)
		}

		if !converted.IsZero() {
			target.Value.SetMapIndex(key, converted)
		}
	}

	return nil
}

func (d *Converter) convertToArray(source, target reflect.Value) error {
	switch kind(source) {
	case reflect.String:
		if ok, err := d.textUnmarshal(source.String(), target); ok {
			if err != nil {
				return rerror(source, target, err)
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
			switch kind(target.Value) {
			case reflect.Array:
				if index >= target.Value.Len() {
					return nil
				}

				return set(target.Value.Index(index), converted)
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

	return set(origin, target)
}

func (d *Converter) convertToBasic(source, target reflect.Value) error {
	if !source.IsValid() {
		source = create(source.Type())
	}

	if target.IsValid() && target.Elem().IsValid() {
		targetElem := refer(target.Elem())
		source = elem(source)

		if err := d.convert(source, targetElem); err != nil {
			return err
		}

		source = targetElem
	}

	return set(target, source)
}

func (d *Converter) textMarshal(source reflect.Value) (string, bool, error) {
	targetType := reflect.TypeOf(new(encoding.TextMarshaler)).Elem()

	for _, variant := range variants(source) {
		if variant.Type().Implements(targetType) {
			var (
				encoder   = variant.Interface().(encoding.TextMarshaler)
				data, err = encoder.MarshalText()
			)

			return string(data), true, err
		}
	}

	return "", false, nil
}

func (d *Converter) textUnmarshal(data string, target reflect.Value) (bool, error) {
	targetType := reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

	for _, variant := range variants(target) {
		if variant.Type().Implements(targetType) {
			var (
				decoder = variant.Interface().(encoding.TextUnmarshaler)
				err     = decoder.UnmarshalText([]byte(data))
			)

			return true, err
		}
	}

	return false, nil
}

func (d *Converter) valueRead(source reflect.Value) (interface{}, bool, error) {
	targetType := reflect.TypeOf(new(driver.Valuer)).Elem()

	for _, variant := range variants(source) {
		if variant.Type().Implements(targetType) {
			var (
				valuer    = variant.Interface().(driver.Valuer)
				data, err = valuer.Value()
			)

			return data, true, err
		}
	}

	return "", false, nil
}

func (d *Converter) valueScan(value interface{}, target reflect.Value) (bool, error) {
	targetType := reflect.TypeOf(new(sql.Scanner)).Elem()

	for _, variant := range variants(target) {
		if variant.Type().Implements(targetType) {
			var (
				scanner = variant.Interface().(sql.Scanner)
				err     = scanner.Scan(value)
			)

			return true, err
		}
	}

	return false, nil
}
