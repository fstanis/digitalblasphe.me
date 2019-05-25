package digitalblasphemy

import (
	"github.com/PuerkitoBio/goquery"
)

func parseIndex(doc *goquery.Document, resolution string) []string {
	var result []string
	doc.Find("a").Each(func(i int, el *goquery.Selection) {
		link, exists := el.Attr("href")
		if !exists {
			return
		}
		if urlRegexpForResolution[resolution].MatchString(link) {
			result = append(result, link)
		}
	})
	return result
}

func parseFreebies(doc *goquery.Document) []string {
	var result []string
	doc.Find("figure > a > img").Each(func(i int, el *goquery.Selection) {
		id, exists := el.Attr("id")
		if !exists {
			return
		}
		result = append(result, id)
	})
	return result
}
