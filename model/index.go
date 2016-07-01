package model

import "github.com/blevesearch/bleve"

// InitIndex initializes the search index at the specified path
func InitIndex(filepath string) (bleve.Index, error) {
	index, err := bleve.Open(filepath)

	// Doesn't yet exist (or error opening) so create a new one
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(filepath, mapping)
		if err != nil {
			return nil, err
		}
	}
	return index, nil
}
