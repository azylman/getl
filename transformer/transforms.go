package transformer

import (
	"github.com/azylman/getl"
)

// TableTransform returns a Table that has applies the given transform function to the output channel.
func TableTransform(input getl.Table, transform func(getl.Row, chan<- getl.Row) error) getl.Table {
	wrappedTransform := func(in <-chan getl.Row, out chan<- getl.Row) error {
		for row := range in {
			if err := transform(row, out); err != nil {
				return err
			}
		}
		return nil
	}
	return getl.Transform(input, wrappedTransform)
}

// Select returns a Table that only has Rows that pass the filter.
func Select(table getl.Table, filter func(getl.Row) (bool, error)) getl.Table {
	return TableTransform(table, func(row getl.Row, out chan<- getl.Row) error {
		pass, err := filter(row)
		if err != nil || !pass {
			return err
		}
		out <- row
		return nil
	})
}

// RowTransform returns a Table that applies a transform function to every row in the input table.
func RowTransform(input getl.Table, transform func(getl.Row) (getl.Row, error)) getl.Table {
	return TableTransform(input, func(in getl.Row, out chan<- getl.Row) error {
		row, err := transform(in)
		if err != nil {
			return err
		}
		out <- row
		return nil
	})
}

// Fieldmap returns a Table that has all the Rows of the input Table with the field mapping applied.
func Fieldmap(table getl.Table, mappings map[string][]string) getl.Table {
	return RowTransform(table, func(row getl.Row) (getl.Row, error) {
		newRow := getl.Row{}
		for key, vals := range mappings {
			for _, val := range vals {
				newRow[val] = row[key]
			}
		}
		return newRow, nil
	})
}

// Valuemap returns a Table that has all the Rows of the input Table with a value mapping applied.
func Valuemap(table getl.Table, mappings map[string]map[interface{}]interface{}) getl.Table {
	return RowTransform(table, func(row getl.Row) (getl.Row, error) {
		newRow := getl.Row{}
		for key, val := range row {
			if mappings[key] == nil || mappings[key][val] == nil {
				newRow[key] = val
				continue
			}
			newRow[key] = mappings[key][val]
		}
		return newRow, nil
	})
}
