package uat

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Browser struct {
	browser *rod.Browser
	page    *rod.Page
	uri     string
	timeout string
}

func NewBrowser(uri string, visible bool, timeout string, args map[string]string) (*Browser, error) {
	if timeout == "" {
		timeout = "3s"
	}
	launcher := launcher.New().Headless(!visible)
	for k, v := range args {
		if v == "" {
			launcher.Set(flags.Flag(k))
		} else {
			launcher.Set(flags.Flag(k), v)
		}
	}
	url, err := launcher.Launch()
	if err != nil {
		return nil, err
	}
	browser := rod.New().ControlURL(url)
	err = browser.Connect()
	if err != nil {
		return nil, err
	}
	page, err := browser.Page(proto.TargetCreateTarget{URL: uri})
	if err != nil {
		return nil, err
	}
	window, err := page.GetWindow()
	if err != nil {
		return nil, err
	}
	err = page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:  *window.Width,
		Height: *window.Height,
	})
	if err != nil {
		return nil, err
	}
	return &Browser{
		browser: browser,
		page:    page,
		uri:     uri,
		timeout: timeout,
	}, nil
}

func (b *Browser) Attr(selector string, scrollable bool, name string) (*string, error) {
	element, err := b.Element(selector, scrollable)
	if err != nil {
		return nil, err
	}
	return element.Attribute(name)
}

func (b *Browser) Click(selector string, scrollable bool, duration ...string) error {
	element, err := b.Element(selector, scrollable)
	if err != nil {
		return err
	}
	err = element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return err
	}
	if len(duration) == 0 {
		return nil
	}
	return b.page.WaitStable(cast.ToDuration(duration[0]))
}

func (b *Browser) Close() error {
	return b.browser.Close()
}

func (b *Browser) Element(selector string, scrollable bool) (*rod.Element, error) {
	var element *rod.Element
	if !scrollable {
		element, _ = b.page.Timeout(cast.ToDuration(b.timeout)).Element(selector)
		if element == nil {
			element, _ = b.page.Timeout(cast.ToDuration(b.timeout)).ElementX(selector)
		}
	} else {
		err := b.scrollToTop()
		if err != nil {
			return nil, err
		}
		var bottom bool
		for {
			element, _ = b.page.Timeout(cast.ToDuration(b.timeout)).Element(selector)
			if element == nil {
				element, _ = b.page.Timeout(cast.ToDuration(b.timeout)).ElementX(selector)
			}
			if element != nil || bottom {
				break
			}
			bottom, err = b.scrollToNext()
			if err != nil {
				return nil, err
			}
		}
	}
	if element == nil {
		return nil, fmt.Errorf("element %s not found", selector)
	}
	return element, nil
}

func (b *Browser) Elements(selector string, scrollable bool) ([]*rod.Element, error) {
	elements, err := b.ElementsFunc([]string{
		selector,
	}, scrollable, nil)
	if err != nil {
		return nil, err
	}
	return lo.Values(elements), nil
}

func (b *Browser) ElementsFunc(selectors []string, scrollable bool, callback func(selector string, element *rod.Element) bool) (map[string]*rod.Element, error) {
	elements := make(map[string]*rod.Element)
	if !scrollable {
		for _, s := range selectors {
			eles, _ := b.page.Timeout(cast.ToDuration(b.timeout)).Elements(s)
			if eles == nil {
				eles, _ = b.page.Timeout(cast.ToDuration(b.timeout)).ElementsX(s)
			}
			for _, e := range eles {
				key, _ := e.HTML()
				if key == "" {
					key = e.String()
				}
				_, ok := elements[key]
				if ok {
					continue
				}
				if callback == nil || callback(s, e) {
					elements[key] = e
				}
			}
		}
	} else {
		err := b.scrollToTop()
		if err != nil {
			return nil, err
		}
		var bottom bool
		for {
			for _, s := range selectors {
				eles, _ := b.page.Timeout(cast.ToDuration(b.timeout)).Elements(s)
				if eles == nil {
					eles, _ = b.page.Timeout(cast.ToDuration(b.timeout)).ElementsX(s)
				}
				for _, e := range eles {
					key, _ := e.HTML()
					if key == "" {
						key = e.String()
					}
					_, ok := elements[key]
					if ok {
						continue
					}
					if callback == nil || callback(s, e) {
						elements[key] = e
					}
				}
			}
			if bottom {
				break
			}
			bottom, err = b.scrollToNext()
			if err != nil {
				return nil, err
			}
		}
	}
	return elements, nil
}

func (b *Browser) Goto(uri string) error {
	url, err := url.JoinPath(b.uri, uri)
	if err != nil {
		return err
	}
	err = b.page.Navigate(url)
	if err != nil {
		return err
	}
	b.page.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)()
	return nil
}

func (b *Browser) Input(selector string, scrollable bool, val string) error {
	element, err := b.Element(selector, scrollable)
	if err != nil {
		return err
	}
	return element.Input(val)
}

func (b *Browser) Text(selector string, scrollable bool) (string, error) {
	element, err := b.Element(selector, scrollable)
	if err != nil {
		return "", err
	}
	return element.Text()
}

func (b *Browser) scrollToNext() (bool, error) {
	height, err := b.page.Eval("() => document.documentElement.clientHeight")
	if err != nil {
		return false, err
	}
	_, err = b.page.Eval(fmt.Sprintf("() => window.scrollBy(0, %d)", height.Value.Int()))
	if err != nil {
		return false, err
	}
	b.page.WaitRequestIdle(300*time.Millisecond, nil, nil, nil)()
	scrollTop, err := b.page.Eval("() => document.documentElement.scrollTop")
	if err != nil {
		return false, err
	}
	scrollHeight, err := b.page.Eval("() => document.body.scrollHeight")
	if err != nil {
		return false, err
	}
	if scrollTop.Value.Int()+height.Value.Int() < scrollHeight.Value.Int() {
		return false, nil
	}
	return true, nil
}

func (b *Browser) scrollToTop() error {
	_, err := b.page.Eval("() => window.scrollTo(0, 0)")
	return err
}
