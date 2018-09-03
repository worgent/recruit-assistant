package utils

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type Spider struct {
	Url string
	Doc *goquery.Document
}

func (s *Spider) getDoc() (*goquery.Document, error) {

	if s.Doc == nil {
		// Request the HTML page.
		res, err := http.Get(s.Url)
		if err != nil {
			log.Println(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		s.Doc = doc
		return doc, err
	}
	return s.Doc, nil
}

func (s *Spider) Text(css string) string {
	doc, err := s.getDoc()
	if err != nil {
		log.Println(err)
	}
	return doc.Find(css).Text()
}

func (s *Spider) Attr(css, attribute string) (string, bool) {
	doc, err := s.getDoc()
	if err != nil {
		log.Println(err)
	}
	return doc.Find(css).Attr(attribute)
}

func (s *Spider) Html(css string) (string, error) {
	doc, err := s.getDoc()
	if err != nil {
		log.Println(err)
	}
	return doc.Find(css).Html()
}
