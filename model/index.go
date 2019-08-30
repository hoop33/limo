package model

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
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

func buildIndexMapping() *mapping.IndexMappingImpl {
	simpleTextFieldMapping := bleve.NewTextFieldMapping()
	simpleTextFieldMapping.Analyzer = simple.Name

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	starMapping := bleve.NewDocumentMapping()
	starMapping.AddFieldMappingsAt("Name", simpleTextFieldMapping)
	starMapping.AddFieldMappingsAt("FullName", simpleTextFieldMapping)
	starMapping.AddFieldMappingsAt("Description", englishTextFieldMapping)
	starMapping.AddFieldMappingsAt("Language", keywordFieldMapping)
	starMapping.AddFieldMappingsAt("Tags.Name", keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Star", starMapping)

	return indexMapping
}
