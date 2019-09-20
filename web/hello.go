package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jackytck/jcconv/detector"
	"github.com/jackytck/jcconv/translator"
)

// Index renders the index page with source represented as page.
func Index(page string, det *detector.Detector, trans2hk, trans2s *translator.Translator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get input from url param if any
		var text, output, error string
		pageWithData := page
		switch r.Method {
		case http.MethodGet:
			texts, ok := r.URL.Query()["text"]
			if ok && len(texts[0]) > 0 {
				text = texts[0]

				// translate in ssr
				isTrad, err := det.IsTraditional(text)
				if err != nil {
					error = err.Error()
					break
				}
				var trans *translator.Translator
				if isTrad {
					trans = trans2s
				} else {
					trans = trans2hk
				}

				output, err = trans.TranslateOne(text)
				if err != nil {
					error = err.Error()
					break
				}
			}
		}
		pageWithData = strings.Replace(page, "{INPUT}", text, 1)
		pageWithData = strings.Replace(pageWithData, "{OUTPUT}", output, 1)
		pageWithData = strings.Replace(pageWithData, "{ERROR}", error, 1)

		fmt.Fprint(w, pageWithData)
	}
}
