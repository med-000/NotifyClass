package scraping

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseClass(html string) *Class {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil
	}

	class := &Class{
		Title: strings.TrimSpace(
			doc.Find("a.course-name").First().Text(),
		),
		Groups: []*Group{},
	}

	doc.Find(".cl-contentsList_folder").Each(func(i int, folder *goquery.Selection) {
		g := &Group{
			Name:   strings.TrimSpace(folder.Find(".panel-title").Text()),
			Events: []*Event{},
		}

		folder.Find(".cl-contentsList_listGroupItem").Each(func(j int, item *goquery.Selection) {
			name := strings.TrimSpace(item.Find(".cm-contentsList_contentName").Text())
			category := strings.TrimSpace(item.Find(".cl-contentsList_categoryLabel").Text())
			date := strings.TrimSpace(item.Find(".cm-contentsList_contentDetailListItemData").Text())
			id := ""

			if a := item.Find(".cl-contentsList_contentDetailListItemData a"); a.Length() > 0 {
				href, _ := a.Attr("href")

				parts := strings.Split(href, "/contents/")
				if len(parts) > 1 {
					id = strings.Split(parts[1], "/")[0]
				}
			}

			e := &Event{
				Id:       id,
				Name:     name,
				Category: category,
				Date:     date,
			}

			g.Events = append(g.Events, e)
		})

		class.Groups = append(class.Groups, g)
	})

	return class
}
