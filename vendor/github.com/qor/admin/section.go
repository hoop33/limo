package admin

import (
	"fmt"
	"strings"

	"github.com/qor/qor/utils"
)

// Section is used to structure forms, it could group your fields into sections, to make your form clean & tidy
//    product.EditAttrs(
//      &admin.Section{
//      	Title: "Basic Information",
//      	Rows: [][]string{
//      		{"Name"},
//      		{"Code", "Price"},
//      	}},
//      &admin.Section{
//      	Title: "Organization",
//      	Rows: [][]string{
//      		{"Category", "Collections", "MadeCountry"},
//      	}},
//      "Description",
//      "ColorVariations",
//    }
type Section struct {
	Resource *Resource
	Title    string
	Rows     [][]string
}

// String stringify section
func (section *Section) String() string {
	return fmt.Sprint(section.Rows)
}

func (res *Resource) generateSections(values ...interface{}) []*Section {
	var sections []*Section
	var hasColumns, excludedColumns []string

	// Reverse values to make the last one as a key one
	// e.g. Name, Code, -Name (`-Name` will get first and will skip `Name`)
	for i := len(values) - 1; i >= 0; i-- {
		value := values[i]
		if section, ok := value.(*Section); ok {
			newSection := uniqueSection(section, &hasColumns)
			if len(newSection.Rows) > 0 {
				sections = append(sections, newSection)
			}
		} else if column, ok := value.(string); ok {
			if strings.HasPrefix(column, "-") {
				excludedColumns = append(excludedColumns, column)
			} else if !isContainsColumn(excludedColumns, column) {
				sections = append(sections, &Section{Rows: [][]string{{column}}})
			}
			hasColumns = append(hasColumns, column)
		} else if row, ok := value.([]string); ok {
			for j := len(row) - 1; j >= 0; j-- {
				column = row[j]
				sections = append(sections, &Section{Rows: [][]string{{column}}})
				hasColumns = append(hasColumns, column)
			}
		} else {
			utils.ExitWithMsg(fmt.Sprintf("Qor Resource: attributes should be Section or String, but it is %+v", value))
		}
	}

	sections = reverseSections(sections)
	for _, section := range sections {
		section.Resource = res
	}
	return sections
}

func uniqueSection(section *Section, hasColumns *[]string) *Section {
	newSection := Section{Title: section.Title}
	var newRows [][]string
	for _, row := range section.Rows {
		var newColumns []string
		for _, column := range row {
			if !isContainsColumn(*hasColumns, column) {
				newColumns = append(newColumns, column)
				*hasColumns = append(*hasColumns, column)
			}
		}
		if len(newColumns) > 0 {
			newRows = append(newRows, newColumns)
		}
	}
	newSection.Rows = newRows
	return &newSection
}

func reverseSections(sections []*Section) []*Section {
	var results []*Section
	for i := 0; i < len(sections); i++ {
		results = append(results, sections[len(sections)-i-1])
	}
	return results
}

func isContainsColumn(hasColumns []string, column string) bool {
	for _, col := range hasColumns {
		if strings.TrimLeft(col, "-") == strings.TrimLeft(column, "-") {
			return true
		}
	}
	return false
}

func containsPositiveValue(values ...interface{}) bool {
	for _, value := range values {
		if _, ok := value.(*Section); ok {
			return true
		} else if column, ok := value.(string); ok {
			if !strings.HasPrefix(column, "-") {
				return true
			}
		} else {
			utils.ExitWithMsg(fmt.Sprintf("Qor Resource: attributes should be Section or String, but it is %+v", value))
		}
	}
	return false
}

// ConvertSectionToMetas convert section to metas
func (res *Resource) ConvertSectionToMetas(sections []*Section) []*Meta {
	var metas []*Meta
	for _, section := range sections {
		for _, row := range section.Rows {
			for _, col := range row {
				meta := res.GetMeta(col)
				if meta != nil {
					metas = append(metas, meta)
				}
			}
		}
	}
	return metas
}

// ConvertSectionToStrings convert section to strings
func (res *Resource) ConvertSectionToStrings(sections []*Section) []string {
	var columns []string
	for _, section := range sections {
		for _, row := range section.Rows {
			for _, col := range row {
				columns = append(columns, col)
			}
		}
	}
	return columns
}

func (res *Resource) setSections(sections *[]*Section, values ...interface{}) {
	if len(values) == 0 {
		if len(*sections) == 0 {
			*sections = res.generateSections(res.allAttrs())
		}
	} else {
		var flattenValues []interface{}

		for _, value := range values {
			if columns, ok := value.([]string); ok {
				for _, column := range columns {
					flattenValues = append(flattenValues, column)
				}
			} else if _sections, ok := value.([]*Section); ok {
				for _, section := range _sections {
					flattenValues = append(flattenValues, section)
				}
			} else if section, ok := value.(*Section); ok {
				flattenValues = append(flattenValues, section)
			} else if column, ok := value.(string); ok {
				flattenValues = append(flattenValues, column)
			} else {
				utils.ExitWithMsg(fmt.Sprintf("Qor Resource: attributes should be Section or String, but it is %+v", value))
			}
		}

		if containsPositiveValue(flattenValues...) {
			*sections = res.generateSections(flattenValues...)
		} else {
			var columns, availbleColumns []string
			for _, value := range flattenValues {
				if column, ok := value.(string); ok {
					columns = append(columns, column)
				}
			}

			for _, column := range res.allAttrs() {
				if !isContainsColumn(columns, column) {
					availbleColumns = append(availbleColumns, column)
				}
			}
			*sections = res.generateSections(availbleColumns)
		}
	}
}
