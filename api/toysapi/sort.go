package toysapi

import (
	"errors"
	"reflect"
	"sort"
	"strings"
)

const (
	orderAsc  string = "asc"
	orderDesc string = "desc"
)

// Sort is a helper struct that groups sorting attributes together with their order
type Sort struct {
	attribute string
	order     string
}

// Sort takes a list of category entries and sorts it due to the given argument list
func doSort(entries []CategoryEntry, args []string) error {

	criteria := make([]Sort, 0, 2)

	for _, arg := range args {

		elemOrder := orderDesc

		els := strings.SplitN(arg, ":", 2)

		attr := transformSortParam(els[0])

		if attr == "" {
			return errors.New("Sorting arguments have invalid format")
		}

		if len(els) > 1 {
			if els[1] != orderAsc && els[1] != orderDesc {
				return errors.New("Invalid sorting criteria")
			}

			elemOrder = els[1]
		}

		criteria = append(criteria, Sort{attr, elemOrder})
	}

	if len(criteria) > 1 && criteria[0].attribute == criteria[1].attribute {
		return errors.New("Invalid sorting criteria")
	}

	sort.Slice(entries, func(i, j int) bool {

		for _, c := range criteria {

			e1 := getFieldString(&entries[i], c.attribute)
			e2 := getFieldString(&entries[j], c.attribute)

			if c.order == orderDesc {
				if e1 > e2 {
					return true
				} else if e1 < e2 {
					return false
				}
			} else {
				if e1 < e2 {
					return true
				} else if e1 > e2 {
					return false
				}
			}
		}

		return true
	})

	return nil
}

func transformSortParam(param string) string {

	if strings.EqualFold(param, "url") {
		return "URL"
	}

	if strings.EqualFold(param, "label") {
		return "Label"
	}

	return ""
}

func getFieldString(c *CategoryEntry, field string) string {
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	return f.String()
}
