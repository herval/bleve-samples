package bleve_samples

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"os"
	"time"
)

type HtmlDocument struct {
	Headers   []string
	Body      string
	CreatedAt time.Time
}

func NewIndex(version string) (bleve.Index, error) {
	var err error
	path := fmt.Sprintf("./test_index" + version + ".bleve")
	mapping := bleve.NewIndexMapping()

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return bleve.New(path, mapping)
	} else {
		return bleve.Open(path)
	}
}
