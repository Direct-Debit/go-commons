package filespec

import (
	"errors"
	"fmt"
	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/Direct-Debit/go-commons/format"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type RecordTag struct {
	Start int
	End   int
	Type  string
}

func ParseRecordTag(f reflect.StructField) (RecordTag, error) {
	positions := strings.Split(f.Tag.Get("pos"), "-")
	if len(positions) < 2 {
		return RecordTag{}, errors.New("pos tag invalid for " + f.Name)
	}

	start, err := strconv.Atoi(positions[0])
	if err != nil {
		return RecordTag{}, err
	}
	end, err := strconv.Atoi(positions[1])
	if err != nil {
		return RecordTag{}, err
	}

	return RecordTag{
		Start: start,
		End:   end,
		Type:  f.Tag.Get("type"),
	}, nil
}

func (r RecordTag) Length() int {
	return r.End - r.Start + 1
}

func parseStruct(field reflect.Value, strVal string, tag RecordTag) (err error) {
	switch field.Type() {
	case reflect.TypeOf(time.Time{}):
		var timeVal time.Time
		switch tag.Length() {
		case 6:
			timeVal, err = time.Parse(format.DateShort6, strVal)
		case 8:
			timeVal, err = time.Parse(format.DateShort8, strVal)
		default:
			err = errors.New(fmt.Sprintf("invalid time length: %d", tag.Length()))
		}
		field.Set(reflect.ValueOf(timeVal))
	default:
		err = errors.New("unsupported struct in parse")
	}
	return err
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
		errlib.FatalError(err, "Invalid tag")

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
				return err
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intVal)
		case reflect.Bool:
			intVal, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return err
			}
			field.SetBool(intVal > 0)
		}

	}

	return nil
}

func structValToStr(val reflect.Value, tag RecordTag) string {
	switch val.Type() {
	case reflect.TypeOf(time.Time{}):
		date := val.Interface().(time.Time)
		switch tag.Length() {
		case 10:
			return date.Format(format.DateShortSlashes)
		case 8:
			return date.Format(format.DateShort8)
		default:
			return date.Format(format.DateShort6)
		}
	}
	errlib.FatalError(
		fmt.Errorf("couldn't convert %v to string", val.Type()),
		"Error with generating file line",
	)
	panic("Fatal error did not fatally exit.")
}

func GenerateLine(source interface{}, builder *strings.Builder) error {
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)

	lastFieldType := sourceType.Field(sourceType.NumField() - 1)
	lastTag, err := ParseRecordTag(lastFieldType)
	errlib.FatalError(err, "Couldn't parse line tag")
	line := make([]rune, lastTag.End+1) // +1 for newline character
	for i := range line {
		line[i] = ' '
	}

	for i := 0; i < sourceValue.NumField(); i++ {
		fieldType := sourceType.Field(i)
		fieldValue := sourceValue.Field(i)

		tag, err := ParseRecordTag(fieldType)
		errlib.FatalError(err, "Couldn't parse line tag")

		var value string
		switch fieldValue.Kind() {
		case reflect.String:
			value = fieldValue.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Struct:
			value = structValToStr(fieldValue, tag)
		}

		switch tag.Type {
		case "N", "D":
			value = fmt.Sprintf("%0*s", tag.Length(), value)
			_, err = strconv.Atoi(value)
			if err != nil && !strings.Contains(value, "TEST") {
				log.Error(err, "Failed for field "+fieldType.Name)
				return err
			}
		case "AN":
			value = strings.ToUpper(value)
			value = fmt.Sprintf("%-*s", tag.Length(), value)
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
