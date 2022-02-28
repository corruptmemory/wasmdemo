package main

import (
	"syscall/js"

	"github.com/Nerzal/tinydom"
)

type item struct {
	name  string
	email string
}

var defaultItems = []item{
	{"Jim", "jim@example.com"},
	{"Therese", "therese@example.com"},
	{"Jimmy", "jimmy@example.com"},
	{"Matthew", "matthew@example.com"},
	{"Brendan", "brendan@example.com"},
	{"Fiona", "fiona@example.com"},
}

type listView struct {
	list         *tinydom.Element
	detail       *tinydom.Element
	listFuncs    []js.Func
	listElements []*tinydom.Element
}

func (l *listView) clear() {
	l.list = nil
	l.detail = nil
	for i := 0; i < len(l.listFuncs); i++ {
		l.listFuncs[i].Release()
	}
	l.listFuncs = nil
	l.listElements = nil
}

type App struct {
	document *tinydom.Document
	window   *tinydom.Window
	history  *tinydom.History
	json     js.Value
	body     *tinydom.Element
	items    []item
	listView listView
}

func (a *App) stringify(in js.Value) string {
	return a.json.Call("stringify", in).String()
}

func (a *App) element(tag, id string) *tinydom.Element {
	e := a.document.CreateElement(tag)
	if id != "" {
		e.SetId(id)
	}
	return e
}

func (a *App) addLink(e *tinydom.Element, text, url string) *tinydom.Element {
	an := a.document.CreateElement("a")
	an.SetInnerHTML(text)
	an.SetAttribute("href", url)
	e.AppendChild(an)
	return an
}

func (a *App) addBr(e *tinydom.Element) *tinydom.Element {
	br := a.document.CreateElement("br")
	e.AppendChild(br)
	return br
}

func (a *App) addH1(e *tinydom.Element, text string) *tinydom.Element {
	h1 := a.document.CreateElement("h1")
	h1.SetInnerHTML(text)
	e.AppendChild(h1)
	return h1
}

func (a *App) getElementsByTag(tag string) []*tinydom.Element {
	v := a.document.Call("getElementsByTagName", tag)
	vl := v.Length()
	if vl > 0 {
		var result []*tinydom.Element
		for i := 0; i < vl; i++ {
			result = append(result, &tinydom.Element{Value: v.Index(i)})
		}
		return result
	}
	return nil
}

func (a *App) clear() {
	a.body.RemoveAllChildNodes()
}

func (a *App) itemListElement(i *item) *tinydom.Element {
	e := a.element("div", "")
	e.SetClass("list-item")
	e.SetInnerHTML(i.name)
	return e
}

func (a *App) updateDetail(i *item) {
	lv := &a.listView
	lv.detail.RemoveAllChildNodes()
	lv.detail.AppendChild(a.element("div", "name").
		SetInnerHTML(i.name))
	lv.detail.AppendChild(a.element("div", "email").
		SetInnerHTML(i.email))
}

func (a *App) setSelected(idx int) {
	lv := &a.listView
	for i := 0; i < len(lv.listElements); i++ {
		e := lv.listElements[i]
		if i == idx {
			e.SetClass("list-item", "selected")
			continue
		}
		e.SetClass("list-item")
	}
}

func (a *App) drawList() {
	a.clear()
	lv := &a.listView
	lv.clear()
	container := a.element("div", "container").SetClass("grid-container")
	a.body.AppendChild(container)
	lv.list = a.element("div", "list")
	container.AppendChild(lv.list)
	lv.detail = a.element("div", "detail")
	container.AppendChild(lv.detail)
	for i := 0; i < len(a.items); i++ {
		itm := &a.items[i]
		ile := a.itemListElement(itm)
		lv.listElements = append(lv.listElements, ile)
		clickFunc := js.FuncOf(func(idx int) func(this js.Value, args []js.Value) interface{} {
			return func(this js.Value, args []js.Value) interface{} {
				a.setSelected(idx)
				a.updateDetail(itm)
				return nil
			}
		}(i))
		ile.AddEventListener("click", clickFunc)
		lv.list.AppendChild(ile)
	}
}

func main() {
	document := tinydom.GetDocument()
	window := tinydom.GetWindow()
	history := window.History()
	app := App{
		document: document,
		window:   window,
		history:  history,
		json:     window.Get("JSON"),
		body:     document.GetElementById("body-component"),
		items:    defaultItems,
	}
	app.drawList()

	wait := make(chan struct{}, 0)
	<-wait
}
