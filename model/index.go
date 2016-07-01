package model

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/language/en"
)

// InitIndex initializes the search index at the specified path
func InitIndex(filepath string) (bleve.Index, error) {
	index, err := bleve.Open(filepath)

	// Doesn't yet exist (or error opening) so create a new one
	if err != nil {
		index, err = bleve.New(filepath, buildIndexMapping())
		if err != nil {
			return nil, err
		}
	}
	return index, nil
}

func buildIndexMapping() *bleve.IndexMapping {
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	starMapping := bleve.NewDocumentMapping()
	starMapping.AddFieldMappingsAt("Name", englishTextFieldMapping)
	starMapping.AddFieldMappingsAt("FullName", englishTextFieldMapping)
	starMapping.AddFieldMappingsAt("Description", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Star", starMapping)

	return indexMapping
}
