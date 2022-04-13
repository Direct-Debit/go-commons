package filespec

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Direct-Debit/go-commons/format"
)

type RecordTag struct {
	Start  int
	End    int
	Type   string
	Format string
}

func ParseRecordTag(f reflect.StructField) (RecordTag, error) {
	positions := strings.Split(f.Tag.Get("pos"), "-")
	if len(positions) < 2 {
		return RecordTag{}, errors.New("pos tag invalid for " + f.Name)
	}

	start, err := strconv.Atoi(positions[0])
	if err != nil {
		return RecordTag{}, errors.Wrapf(err, "could not parse field start index")
	}
	end, err := strconv.Atoi(positions[1])
	if err != nil {
		return RecordTag{}, errors.Wrapf(err, "could not parse field end index")
	}

	return RecordTag{
		Start:  start,
		End:    end,
		Type:   f.Tag.Get("type"),
		Format: f.Tag.Get("format"),
	}, nil
}

func (r RecordTag) Length() int {
	return r.End - r.Start + 1
}

func parseStruct(field reflect.Value, strVal string, tag RecordTag) error {
	switch field.Type() {
	case reflect.TypeOf(time.Time{}):
		timeFormat := tag.Format
		if len(timeFormat) == 0 {
			switch tag.Length() {
			case 6:
				timeFormat = format.DateShort6
			case 8:
				timeFormat = format.DateShort8
			case 4:
				timeFormat = format.MMYY
			default:
				return fmt.Errorf("invalid time length: %d", tag.Length())
			}
		}
		timeVal, err := time.Parse(timeFormat, strVal)
		if err != nil {
			return errors.Wrapf(err, "could not parse time %s with %s", strVal, timeFormat)
		}
		field.Set(reflect.ValueOf(timeVal))
	default:
		return errors.New("unsupported struct in parse")
	}
	return nil
}

func ParseRecord(line string, target interface{}) error {
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return errors.New("target is not a pointer")
	}

	targetType := reflect.TypeOf(target).Elem()
	targetValue := reflect.ValueOf(target).Elem()

	for i := 0; i < targetValue.NumField(); i++ {
		field := targetValue.Field(i)
		fieldType := targetType.Field(i)

		tag, err := ParseRecordTag(fieldType)
		if err != nil {
			return errors.Wrapf(err, "invalid tag on %s", fieldType.Name)
		}

		strVal := line[tag.Start-1 : tag.End]

		switch field.Kind() {
		case reflect.String:
			val := strVal
			if tag.Type == "N" {
				val = strings.TrimSpace(val)
				val = strings.TrimLeft(val, "0")
			}
			field.SetString(val)
		case reflect.Struct:
			err = parseStruct(field, strVal, tag)
			if err != nil {
				return errors.Wrapf(err, "could not parse struct for %s", fieldType.Name)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return errors.Wrapf(err, "could not parse int for %s", fieldType.Name)
			}
			field.SetInt(intVal)
		case reflect.Bool:
			switch tag.Type {
			case "N":
				intVal, err := strconv.ParseInt(strVal, 10, 64)
				if err != nil {
					return errors.Wrapf(err, "could not parse bool for %s", fieldType.Name)
				}
				field.SetBool(intVal > 0)
			case "AN", "A":
				field.SetBool(true)
				val := strings.ToLower(strings.TrimSpace(strVal))
				for _, s := range []string{
					"n", "no",
					"f", "false",
					"0",
				} {
					if s == val {
						field.SetBool(false)
					}
				}
			}
		}
	}

	return nil
}

func structValToStr(val reflect.Value, tag RecordTag) (string, error) {
	switch val.Type() {
	case reflect.TypeOf(time.Time{}):
		date := val.Interface().(time.Time)
		if tag.Format != "" {
			return date.Format(tag.Format), nil
		} else {
			switch tag.Length() {
			case 10:
				return date.Format(format.DateShortSlashes), nil
			case 8:
				return date.Format(format.DateShort8), nil
			default:
				return date.Format(format.DateShort6), nil
			}
		}
	}
	return "", fmt.Errorf("couldn't convert %v to string", val.Type())
}

func GenerateLine(source interface{}, builder *strings.Builder) error {
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)

	lastFieldType := sourceType.Field(sourceType.NumField() - 1)
	lastTag, err := ParseRecordTag(lastFieldType)
	if err != nil {
		return errors.Wrapf(err, "invalid tag on %s", lastFieldType.Name)
	}

	line := make([]rune, lastTag.End+1) // +1 for newline character
	for i := range line {
		line[i] = ' '
	}

	for i := 0; i < sourceValue.NumField(); i++ {
		fieldType := sourceType.Field(i)
		fieldValue := sourceValue.Field(i)

		tag, err := ParseRecordTag(fieldType)
		if err != nil {
			return errors.Wrapf(err, "invalid tag on %s", fieldType.Name)
		}

		var value string
		switch fieldValue.Kind() {
		case reflect.String:
			value = fieldValue.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Struct:
			value, err = structValToStr(fieldValue, tag)
			if err != nil {
				return errors.Wrapf(err, "could not parse struct type to string")
			}
		}

		switch tag.Type {
		case "N":
			value = fmt.Sprintf("%0*s", tag.Length(), value)
			_, err = strconv.Atoi(value)
			if err != nil && !strings.Contains(value, "TEST") {
				return errors.Wrapf(err, "could not convert integer value from %s", fieldType.Name)
			}
		case "C":
			value = fmt.Sprintf("%0*s", tag.Length(), value)
			_, err = strconv.ParseFloat(value, 64)
			if err != nil && !strings.Contains(value, "TEST") {
				return errors.Wrapf(err, "could not convert float value from %s", fieldType.Name)
			}
		case "A", "AN":
			value = strings.ToUpper(value)
			switch tag.Format {
			case "align-right":
				value = fmt.Sprintf("%*s", tag.Length(), value)
			default:
				value = fmt.Sprintf("%-*s", tag.Length(), value) // Default to left align
			}
		}

		for idx, c := range value {
			lineIdx := tag.Start - 1 + idx
			line[lineIdx] = c
		}
	}

	line[len(line)-1] = '\n'
	builder.WriteString(string(line))
	return nil
}
