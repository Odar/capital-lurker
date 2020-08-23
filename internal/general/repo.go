package general

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/Odar/capital-lurker/pkg/api"
)

func ApplyFilter(tableName string, filter *api.Filter, builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter == nil {
		return builder
	}
	if filter.ID != nil { //add fix to query: id starts from 1
		builder = builder.Where(tableName+".id = ?", *filter.ID)
	}
	if filter.Name != nil {
		builder = builder.Where(tableName+".name LIKE ?", "%"+*filter.Name+"%")
	}
	if filter.OnMainPage != nil { //how to parse blanks?
		builder = builder.Where(tableName+".on_main_page = ?", *filter.OnMainPage)
	}
	if filter.InFilter != nil {
		builder = builder.Where(tableName+".in_filter = ?", *filter.InFilter)
	}
	if filter.AddedAtRange != nil {
		builder = builder.Where(tableName+".added_at >= ? AND "+
			tableName+".added_at < ?", filter.AddedAtRange.From, filter.AddedAtRange.To)
	}
	if filter.UpdatedAtRange != nil {
		builder = builder.Where(tableName+".updated_at >= ? AND "+
			tableName+".updated_at < ?", filter.UpdatedAtRange.From, filter.UpdatedAtRange.To)
	}
	if filter.Position != nil {
		builder = builder.Where(tableName+".position = ?", *filter.Position)
	}
	if filter.Img != nil {
		builder = builder.Where(tableName+".img LIKE ?", "%"+*filter.Img+"%")
	}
	return builder
}

func ApplySortByParameter(sortBy string) string {
	var columnNames = map[string]bool{
		"id":           true,
		"name":         true,
		"on_main_page": true,
		"in_filter":    true,
		"added_at":     true,
		"updated_at":   true,
		"position":     true,
		"img":          true,
	}
	var orderByKeywords = map[string]bool{
		"DESC": true,
		"ASC":  true,
	}
	words := strings.Split(sortBy, " ")
	_, foundColumnName := columnNames[words[0]]
	if len(words) == 1 {
		if !foundColumnName {
			sortBy = "id DESC"
		}
	} else {
		_, foundOrderByKeyword := orderByKeywords[words[1]]
		if !foundColumnName || sortBy == "" || len(words) > 2 || !foundOrderByKeyword {
			sortBy = "id DESC"
		}
	}

	return sortBy
}
