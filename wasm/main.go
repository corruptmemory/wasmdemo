package main

import (
	"syscall/js"

	"github.com/Nerzal/tinydom"
	"github.com/Nerzal/tinydom/elements/form"
	"github.com/Nerzal/tinydom/elements/input"
)

type App struct {
	document *tinydom.Document
	window   *tinydom.Window
	history  *tinydom.History
	json     js.Value
}

func (a *App) stringify(in js.Value) string {
	return a.json.Call("stringify", in).String()
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

func main() {
	document := tinydom.GetDocument()
	window := tinydom.GetWindow()
	history := window.History()
	app := App{
		document: document,
		window:   window,
		history:  history,
		json:     window.Get("JSON"),
	}
	println(window.Location().Href())

	body := document.GetElementById("body-component")
	app.addH1(body, "Welcome to tinydom - Hello TinyWorld <3")
	app.addH1(body, "Yes! I do compile with TinyGo!")
	app.addBr(body)
	app.addBr(body)
	app.addLink(body, "Link 1", "address1")
	app.addBr(body)
	app.addLink(body, "Link 2", "address2")
	app.addBr(body)
	myForm := form.New()
	label := document.CreateElement("label")
	label.SetInnerHTML("Name:")
	textInput := input.NewTextInput()

	passwordLabel := document.CreateElement("label")
	passwordLabel.SetInnerHTML("Password:")
	passwordInput := input.New(input.PasswordInput)

	submitInput := input.New(input.SubmitInput)

	err := myForm.Append(label, textInput.Element, passwordLabel, passwordInput.Element, submitInput.Element)
	if err != nil {
		println(err.Error())
	}

	body.AppendChild(myForm.Element)
	anchors := app.getElementsByTag("a")
	for i, a := range anchors {
		ok, href := a.GetAttribute("href")
		println("anchor[", i, "]: ", ok, " - ", href)
	}
	for _, a := range anchors {
		a.AddEventListener("click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			element := tinydom.Element{Value: this}
			ok, href := element.GetAttribute("href")
			println("clicked: ", ok, " - ", href)
			event := tinydom.Event{Value: args[0]}
			window.PushState(nil, "", href)
			event.PreventDefault()
			event.StopPropagation()
			return nil
		}))
	}

	wait := make(chan struct{}, 0)
	<-wait
}
