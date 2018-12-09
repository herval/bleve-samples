package bleve_samples

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/simple"
	"os"
	"time"
)

type HtmlDocument struct {
	Headers   []string
	Body      string
	CreatedAt time.Time
	Source    string
	Type      string
}

func NewIndex(version string) (bleve.Index, error) {
	var err error
	path := fmt.Sprintf("./test_index" + version + ".bleve")
	mapping := bleve.NewIndexMapping()

	mapping.DefaultMapping = bleve.NewDocumentStaticMapping() // explicit mapping required for all fields

	htmlDocMapping := bleve.NewDocumentMapping()

	// Source field should be searchable with lowercase and accept full matches only
	lowerCase := bleve.NewTextFieldMapping()
	lowerCase.Analyzer = simple.Name
	htmlDocMapping.AddFieldMappingsAt("Source", lowerCase)

	// Body field should strip html tags
	htmlStripped := bleve.NewTextFieldMapping()
	//if err = mapping.AddCustomAnalyzer("htmlAnalyser", map[string]interface{}{
	//	"type":          custom.Name,
	//	"char_filters":  []string{html.Name},
	//	"tokenizer":     web.Name,
	//	"token_filters": []string{lowercase.Name, camelcase.Name},
	//}); err != nil {
	//	return nil, err
	//}
	//htmlStripped.Analyzer = "htmlAnalyser"
	htmlDocMapping.AddFieldMappingsAt("Body", htmlStripped)

	// CreatedAt field should be searchable by time ranges
	dateTime := bleve.NewDateTimeFieldMapping()
	htmlDocMapping.AddFieldMappingsAt("CreatedAt", dateTime)

	// the "type" field specifies the doc type
	mapping.AddDocumentMapping("html", htmlDocMapping)
	mapping.TypeField = "Type"

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return bleve.New(path, mapping)
	} else {
		return bleve.Open(path)
	}
}
