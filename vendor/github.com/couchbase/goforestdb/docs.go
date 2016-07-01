/*
ForestDB bindings for go

Sample usage (without proper error handling):

	// Open a database
	db, _ := Open("test", nil)

	// Close it properly when we're done
	defer db.Close()

	// Store the document
	doc, _ := NewDoc([]byte("key"), nil, []byte("value"))
	defer doc.Close()
	db.Set(doc)

	// Lookup the document
	doc2, _ := NewDoc([]byte("key"), nil, nil)
	defer doc2.Close()
	db.Get(doc2)

	// Delete the document
	doc3, _ := NewDoc([]byte("key"), nil, nil)
	defer doc3.Close()
	db.Delete(doc3)
*/
package forestdb
