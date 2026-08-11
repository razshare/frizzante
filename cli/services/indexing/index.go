package indexing

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"sync"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gocolly/colly/v2"
	"github.com/razshare/frizzante/v2/tui/messages"
)

func Index(options IndexOptions) (pages map[string]IndexedPage, err error) {
	if options.Depth < 0 {
		return
	}
	var addressUrl *url.URL
	if addressUrl, err = url.Parse(options.Address); err != nil {
		return
	}
	if slices.Contains(Banned, addressUrl.Hostname()) {
		return
	}
	var value any
	var mutex *sync.Mutex
	if value = options.Context.Value("mutex"); value != nil {
		mutex = value.(*sync.Mutex)
	} else {
		mutex = &sync.Mutex{}
		options.Context = context.WithValue(options.Context, "mutex", mutex)
	}
	var tracked map[string]IndexedPage
	if value = options.Context.Value("tracked"); value != nil {
		tracked = *value.(*map[string]IndexedPage)
	} else {
		tracked = make(map[string]IndexedPage)
		options.Context = context.WithValue(options.Context, "tracked", &tracked)
	}
	context.WithValue(options.Context, "mutex", mutex)
	mutex.Lock()
	if _, exists := tracked[options.Address]; exists {
		mutex.Unlock()
		return
	}
	mutex.Unlock()
	var current int
	var maximum int
	mutex.Lock()
	if tracked == nil {
		tracked = make(map[string]IndexedPage)
	}
	mutex.Unlock()
	var title string
	var body string
	addresses := make([]string, 0)
	collector := colly.NewCollector()
	var collectorErrors []error
	collector.OnHTML("a[href]", func(element *colly.HTMLElement) {
		href := element.Attr("href")
		if strings.TrimSpace(href) == "" {
			return
		}
		if !strings.HasPrefix(href, "http://") && !strings.HasPrefix(href, "https://") {
			href = fmt.Sprintf("%s://%s%s", addressUrl.Scheme, addressUrl.Host, href)
		}
		if options.StickToHost {
			var hrefUrl *url.URL
			if hrefUrl, err = url.Parse(href); err != nil {
				return
			}
			if hrefUrl.Host != addressUrl.Host {
				return
			}
		}
		maximum += 1
		addresses = append(addresses, href)
	})
	collector.OnHTML("title", func(element *colly.HTMLElement) {
		maximum += 1
		current += 1
		title = element.Text
	})
	collector.OnHTML("body", func(element *colly.HTMLElement) {
		maximum += 1
		current += 1
		var html string
		var htmlError error
		if html, htmlError = element.DOM.Html(); err != nil {
			collectorErrors = append(collectorErrors, htmlError)
			return
		}
		if body, err = htmltomarkdown.ConvertString(html); err != nil {
			collectorErrors = append(collectorErrors, htmlError)
			return
		}
	})
	messages.Infof("indexing %s", strings.TrimSpace(options.Address))
	if err = collector.Visit(options.Address); err != nil {
		return
	}
	if len(collectorErrors) > 0 {
		err = errors.Join(collectorErrors...)
		return
	}
	mutex.Lock()
	tracked[options.Address] = IndexedPage{
		Title: title,
		Body:  body,
	}
	mutex.Unlock()
	var mutexLocal sync.Mutex
	var group sync.WaitGroup
	var indexErrors []error
	for _, addressLocal := range addresses {
		group.Go(func() {
			if _, indexError := Index(IndexOptions{
				Context:     options.Context,
				Address:     addressLocal,
				Depth:       options.Depth - 1,
				StickToHost: options.StickToHost,
			}); err != nil {
				mutexLocal.Lock()
				indexErrors = append(indexErrors, indexError)
				mutexLocal.Unlock()
				return
			}
			if options.OnProgress != nil {
				group.Go(func() {
					mutexLocal.Lock()
					defer mutexLocal.Unlock()
					current++
					options.OnProgress(current, maximum)
				})
			}
		})
	}
	group.Wait()
	if len(indexErrors) > 0 {
		err = errors.Join(indexErrors...)
	}
	pages = tracked
	return
}
