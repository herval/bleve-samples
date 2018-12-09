package bleve_samples

import (
	"github.com/blevesearch/bleve"
	"log"
	"strconv"
	"testing"
	"time"
)

var index bleve.Index

func TestMain(m *testing.M) {
	i, err := NewIndex(strconv.Itoa(int(time.Now().Unix())))
	if err != nil {
		log.Fatal("Cannot setup index", err)
		return
	}

	// sample content
	if err := i.Index("content 1",
		&HtmlDocument{
			Headers:   []string{"header1", "header2"},
			Body:      "<b>foo bar</b> baz/boz",
			CreatedAt: time.Date(2018, 1, 1, 1, 1, 0, 0, time.UTC),
			Source:    "trusted_Site",
			Type:      "html",
		}); err != nil {
		log.Fatal(err)
		return
	}

	// this should not be searchable because the type is not valid
	if err := i.Index("content 2",
		&HtmlDocument{
			Headers:   []string{"header1", "header2"},
			Body:      "<b>foo bar</b> baz/boz",
			CreatedAt: time.Date(2018, 1, 1, 1, 1, 0, 0, time.UTC),
			Source:    "untrusted_site",
			Type:      "invalid",
		}); err != nil {
		log.Fatal(err)
		return
	}

	index = i
	m.Run()
	_ = index.Close()
}

func TestExactMatchUnsearchableDoc(t *testing.T) {
	q := bleve.NewMatchQuery("untrusted_site")
	q.SetField("Source")

	res, err := index.Search(&bleve.SearchRequest{
		Query: q,
	});
	if err != nil {
		t.Fatal(err)
	}

	if res.Total != 0 {
		t.Fatal("expected to find no document, found: ", res)
	}
}

func TestExactMatch(t *testing.T) {
	// match Source because it's stored as lowercase and the search is normalized to lowercase too
	q := bleve.NewMatchQuery("Trusted_site")
	q.SetField("Source")

	res, err := index.Search(&bleve.SearchRequest{
		Query: q,
	});
	if err != nil {
		t.Fatal(err)
	}

	if res.Total != 1 {
		t.Fatal("expected to find one document, found: ", res)
	}
}

func TestNoExactMatch(t *testing.T) {
	q := bleve.NewMatchQuery("trusted_sit")
	q.SetField("Source")

	res, err := index.Search(&bleve.SearchRequest{
		Query: q,
	});
	if err != nil {
		t.Fatal(err)
	}

	if res.Total != 1 {
		t.Fatal("expected to find no document, found: ", res)
	}
}

func TestHtmlTagsStripped(t *testing.T) {
	res, err := index.Search(&bleve.SearchRequest{
		Query: bleve.NewPrefixQuery("<b>"),
	});
	if err != nil {
		t.Fatal(err)
	}

	if res.Total != 0 {
		t.Fatal("expected to find no document, found: ", res)
	}
}

func TestKeywordsTags(t *testing.T) {
	res, err := index.Search(&bleve.SearchRequest{
		Query: bleve.NewMatchQuery("baz boz"),
	});
	if err != nil {
		t.Fatal(err)
	}

	if res.Total != 1 {
		t.Fatal("expected to find one document, found: ", res)
	}
}

func TestHtmlTags(t *testing.T) {
	res, err := index.Search(&bleve.SearchRequest{
		Query: bleve.NewMatchQuery("foo bar baz"),
	});
	if err != nil {
		t.Fatal(err)
	}

	if res.Total != 1 {
		t.Fatal("expected to find one document, found: ", res)
	}
}

func TestTimeSearch(t *testing.T) {
	true := true
	t0 := time.Date(2018, 1, 1, 1, 0, 0, 0, time.UTC)
	t1 := time.Date(2018, 1, 1, 1, 2, 0, 0, time.UTC)
	res, err := index.Search(&bleve.SearchRequest{
		Query: bleve.NewDateRangeInclusiveQuery(t0, t1, &true, &true),
	});
	if err != nil {
		t.Fatal(err)
	}

	if res.Total != 1 {
		t.Fatal("expected to find one document, found: ", res)
	}
}
