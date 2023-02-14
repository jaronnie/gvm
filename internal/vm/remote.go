package vm

import (
	"path/filepath"

	colly "github.com/gocolly/colly/v2"

	"github.com/jaronnie/gvm/internal/global"
)

type RemoteVM struct {
	Registry string
}

func NewRemoteVM(vm *RemoteVM) Interface {
	return vm
}

func (o RemoteVM) List() ([]string, error) {
	c := colly.NewCollector(
		colly.Async(true),
		colly.CacheDir(filepath.Join(global.GVM_CONFIG_DIR, ".cache")),
	)

	var all []string

	c.OnHTML("h2#stable", func(e *colly.HTMLElement) {
		baseDom := e.DOM.Next()
		for {
			val, exists := baseDom.Attr("id")
			if !exists || val == "archive" {
				break
			}
			all = append(all, val)
			baseDom = baseDom.Next()
		}
	})

	c.OnHTML("div.expanded", func(e *colly.HTMLElement) {
		e.ForEach("div.toggle", func(i int, element *colly.HTMLElement) {
			all = append(all, element.Attr("id"))
		})
	})

	c.OnRequest(func(r *colly.Request) {})

	c.OnResponse(func(r *colly.Response) {})

	err := c.Visit(o.Registry)
	if err != nil {
		return nil, err
	}

	c.Wait()

	return all, nil
}
