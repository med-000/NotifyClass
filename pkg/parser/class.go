package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseClass course.phpのhtml専用
func (p *Parser) ParseClass(html string) *Class {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		p.log.Error.Printf("Cannt Read Html \n Error Detail:%s", err)
		return nil
	}

	//class宣言とpageからClassName取得
	class := &Class{
		Title: strings.TrimSpace(
			doc.Find("a.course-name").First().Text(),
		),
		Groups: []*Group{},
	}

	//folderからgroup取得
	doc.Find(".cl-contentsList_folder").Each(func(i int, folder *goquery.Selection) {
		//groupNameを打ち込む
		g := &Group{
			Name:   strings.TrimSpace(folder.Find(".panel-title").Text()),
			Events: []*Event{},
		}

		//group内のeventを入れる
		folder.Find(".cl-contentsList_listGroupItem").Each(func(j int, item *goquery.Selection) {
			name := strings.TrimSpace(item.Find(".cm-contentsList_contentName").Text())
			category := strings.TrimSpace(item.Find(".cl-contentsList_categoryLabel").Text())
			date := strings.TrimSpace(item.Find(".cm-contentsList_contentDetailListItemData").Text())
			var id string
			var fullURL string

			if a := item.Find(".cl-contentsList_contentDetailListItemData a"); a.Length() > 0 {
				href, _ := a.Attr("href")

				//contentのidを取得
				parts := strings.Split(href, "/contents/")
				if len(parts) > 1 {
					id = strings.Split(parts[1], "/")[0]
					if id == "" {
						p.log.Error.Printf("Id is nil")
					}
				}
			}
			if a := item.Find(".cl-contentsList_contentInfo a"); a.Length() > 0 {
				link, _ := a.Attr("href")
				fullURL = "https://els.sa.dendai.ac.jp" + link
				if link == "" {
					p.log.Error.Printf("Link is nil")
				}
				p.log.Error.Printf("Debug watch %s", fullURL)
			}

			e := &Event{
				ExternalId: id,
				Name:       name,
				Category:   category,
				URL:        fullURL,
				Date:       date,
			}

			g.Events = append(g.Events, e)
		})

		class.Groups = append(class.Groups, g)
	})

	return class
}
