package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Displayable interface {
	Columns() []string
	Values() map[string]interface{}
}

type Column struct {
	Index int
	Name  string
	Width int
}

type Table struct {
	Columns []*Column
	Rows    [][]string
}

func (c *Column) format() string {
	return fmt.Sprintf("%%%ds", -c.Width)
}

func (t *Table) FindColumn(name string) *Column {
	for _, col := range t.Columns {
		if col.Name == name {
			return col
		}
	}
	return nil
}

func (t *Table) insertColumns(cols []string) {
	for idx, col := range cols {
		t.Columns = append(t.Columns, &Column{
			Index: idx,
			Name:  col,
			Width: len(col),
		})
	}
}

func (t *Table) insertRow(data map[string]interface{}) {
	row := make([]string, len(t.Columns))

	for key, val := range data {
		col := t.FindColumn(key)
		if col == nil {
			continue
		}

		str := fmt.Sprintf("%+v", val)
		row[col.Index] = str

		if len(str) > col.Width {
			col.Width = len(str)
		}
	}

	t.Rows = append(t.Rows, row)
}

func (t *Table) Format(writer io.Writer, separator string, pretty bool) error {
	format := "%s"
	for idx, col := range t.Columns {
		if pretty {
			format = col.format()
		}

		_, err := fmt.Fprintf(writer, format, strings.ToUpper(col.Name))
		if err != nil {
			return err
		}

		if (idx + 1) < len(t.Columns) {
			_, err := fmt.Fprintf(writer, separator)
			if err != nil {
				return err
			}
		}
	}

	_, err := fmt.Fprintln(writer)
	if err != nil {
		return err
	}

	for _, row := range t.Rows {
		for idx, val := range row {
			if pretty {
				format = t.Columns[idx].format()
			}

			if strings.Contains(val, separator) {
				val = fmt.Sprintf("%q", val)
			}

			_, err := fmt.Fprintf(writer, format, val)
			if err != nil {
				return err
			}

			if (idx + 1) < len(row) {
				_, err := fmt.Fprint(writer, separator)
				if err != nil {
					return err
				}
			}
		}

		_, err := fmt.Fprintln(writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Table) insertMap(value reflect.Value) error {
	if t.Columns == nil {
		var cols []string
		for _, key := range value.MapKeys() {
			cols = append(cols, fmt.Sprintf("%v", key))
		}
		t.insertColumns(cols)
	}

	row := make(map[string]interface{})

	iter := value.MapRange()
	for iter.Next() {
		row[fmt.Sprintf("%v", iter.Key().Interface())] = iter.Value().Interface()
	}

	t.insertRow(row)
	return nil
}

func (t *Table) insertStruct(value reflect.Value) error {
	if !value.Type().AssignableTo(reflect.TypeOf((*Displayable)(nil)).Elem()) {
		return fmt.Errorf("unable to serialize non `Displayable` struct of type %q", value.Type().String())
	}

	if t.Columns == nil {
		columnsFunc := value.MethodByName("Columns")
		columns := columnsFunc.Call([]reflect.Value{})[0]

		var cols []string
		for i := 0; i < columns.Len(); i++ {
			cols = append(cols, columns.Index(i).String())
		}
		t.insertColumns(cols)
	}

	valuesFunc := value.MethodByName("Values")
	values := valuesFunc.Call([]reflect.Value{})[0]
	return t.insertMap(values)
}

func (t *Table) insertValue(value reflect.Value) error {
	switch value.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			err := t.insertValue(value.Index(i))
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		return t.insertMap(value)
	case reflect.Ptr:
		fallthrough
	case reflect.Struct:
		return t.insertStruct(value)
	}

	return fmt.Errorf("unable to serialize value of type %q (%q)", value.Type().String(), value.Kind().String())
}

func (t *Table) Insert(val interface{}) error {
	return t.insertValue(reflect.ValueOf(val))
}
