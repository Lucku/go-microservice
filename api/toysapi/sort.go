package toysapi

import (
	"errors"
	"net/http"
	"reflect"
	"sort"
	"strings"
)

const (
	orderAsc  string = "asc"
	orderDesc string = "desc"
)

type Sort struct {
	attribute string
	order     string
}

func handleSortParam(input []CategoryEntry, w http.ResponseWriter, r *http.Request) error {

	sortParam := r.URL.Query().Get("sort")

	if sortParam == "" {
		return nil
	}

	sortArgs := strings.Split(sortParam, ",")

	if len(sortArgs) > 2 {
		return errors.New("Invalid numbers of sorting arguments")
	}

	criteria := make([]Sort, 0, 2)

	for _, arg := range sortArgs {

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

	doSort(input, criteria)

	return nil
}

func doSort(input []CategoryEntry, criteria []Sort) {

	sort.Slice(input, func(i, j int) bool {

		for _, c := range criteria {

			e1 := getFieldString(&input[i], c.attribute)
			e2 := getFieldString(&input[j], c.attribute)

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
