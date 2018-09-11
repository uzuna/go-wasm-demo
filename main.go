package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"syscall/js"
	"time"

	"github.com/fogleman/gg"
	"github.com/uzuna/go-wasm/goth"
)

func main() {
	fmt.Println("Hello wasm")
	content := js.Global().Get("document").Call("getElementById", "content")
	cnode := goth.NewNode(content)
	{
		var cb js.Callback

		cb = js.NewCallback(func(args []js.Value) {
			fmt.Println("button clicked")

			// callbackを閉じる。 once
			cb.Release()
		})
		button := goth.CreateElement("button").
			SetAttribute("id", "myButton").
			Set("innerHTML", "Go Callback Test").
			AddEventListener("click", cb)
		goth.AppendChild(cnode, button)
		// js.Global().Get("document").Call("getElementById", "myButton").Call("addEventListener", "click", cb)
	}

	{
		var cb js.Callback

		cb = js.NewCallback(func(args []js.Value) {
			/*
			 以下はjavascriptの操作

			 ```js
			 const content = document.getElementById('id');
			 let div = document.createElement('div');
			 div.textContent = 'new row'
			 content.appendChild(div)
			 ````
			*/

			// get content

			// create div
			div := goth.CreateElement("div")
			div.Set("textContent", "new row")

			// append child
			goth.AppendChild(cnode, div)

			// Releaseしていないのでclickのたびに実行する
			// cb.Release()
		})
		button := goth.CreateElement("button").
			SetAttribute("id", "createDiv").
			Set("innerHTML", "Go Create DOM").
			AddEventListener("click", cb)
		goth.AppendChild(cnode, button)
		// js.Global().Get("document").Call("getElementById", "createDiv").Call("addEventListener", "click", cb)
	}

	// fmt.Println(js.Global().Get("document").Call("getElementById", "myButton").Get("InnerHTML").String())

	// alert := js.Global().Get("alert")
	// alert.Invoke("Alert wasm")
	counter := 0
	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(fmt.Sprintf("Hello wasm inner loop At: %d", counter))
			if counter > 10 {
				fmt.Println("break inner loop")
				break
			}
			counter++
		}
	}()

	// webcomponent
	// reference https://matthewphillips.info/programming/wasm-golang-ce.html
	{
		init := js.NewCallback(func(i []js.Value) {
			el := goth.NewNode(i[0])
			// el.AttachShadow()

			// shadowroot
			// shr := el.CreateShadowRoot()
			// shr.Set("innerHTML", `<style>h3{ color: red; }</style><h3>Shadow DOM</h3>`)

			// layout
			form := goth.CreateElement("form")
			el.AppendChild(form)
			fieldset := goth.CreateElement("fieldset") //style要素
			form.AppendChild(fieldset)

			// view
			label := goth.CreateElement("label").
				SetAttribute("for", "title").
				Set("innerHTML", "UserName")
			fieldset.AppendChild(label)

			inputText := goth.CreateElement("input").
				SetAttribute("type", "text")
			fieldset.AppendChild(inputText)
			// view
			labelPass := goth.CreateElement("label").
				SetAttribute("for", "title").
				Set("innerHTML", "Password")
			fieldset.AppendChild(labelPass)

			inputTextPass := goth.CreateElement("input").
				SetAttribute("type", "password")
			fieldset.AppendChild(inputTextPass)

			// event

			cb := js.NewEventCallback(js.PreventDefault, func(event js.Value) {
				uname := inputText.Get("value").String()
				pass := inputTextPass.Get("value")
				fmt.Println("CallBack", uname, pass)
			})

			submit := goth.CreateElement("button").
				Set("innerHTML", "Submit")
			fieldset.AppendChild(submit)

			form.AddEventListener("submit", cb)

		})

		js.Global().Call("makeComponent", "hello-world", init)

		hw := goth.CreateElement("hello-world")
		goth.AppendChild(cnode, hw)
	}

	// web component with html template
	// js libと同じようにHTML template 解析でイベントを受ける
	// structにmapできるならよいか?
	{

		html := `<form>
			<fieldset>
			<label>Username</label>
			<input type="text" name="username" value="guest">
			<label>Password</label>
			<input type="password" name="password" value="pass">
			<button>Submit</button>
			</fieldset>
		</form>`
		init := js.NewCallback(func(i []js.Value) {
			el := goth.NewNode(i[0])
			el.Set("innerHTML", html)

			//

			fmt.Println(el.Get("children"))
			fmt.Println(el.Get("children").Length())
			form := el.Get("children").Index(0)

			fmt.Println(form.Get("children"))
			fmt.Println(form.Get("children").Length())
			fieldset := form.Get("children").Index(0)

			fmt.Println(fieldset.Get("children"))
			fmt.Println(fieldset.Get("children").Length())

			// match template -> TypeError: Right-hand side of 'instanceof' is not callable
			// instanceInput := js.Global().Get("document").Call("createElement", "input")
			// fmt.Println(instanceInput)
			// fmt.Println(instanceInput.Type())

			var inputText, inputTextPass js.Value

			for i := 0; i < fieldset.Get("children").Length(); i++ {
				d := fieldset.Get("children").Index(i)
				// d.InstanceOf(instanceInput) // TypeError
				fmt.Println(i, d, d.Get("name"), d.Get("tagName"))
				if d.Get("tagName").String() == "INPUT" && d.Get("name").String() == "username" {
					inputText = d
				} else if d.Get("tagName").String() == "INPUT" && d.Get("name").String() == "password" {
					inputTextPass = d
				}
			}

			cb := js.NewEventCallback(js.PreventDefault, func(event js.Value) {
				uname := inputText.Get("value").String()
				pass := inputTextPass.Get("value")
				fmt.Println("CallBack", uname, pass)
			})

			elform := goth.NewNode(form)
			elform.AddEventListener("submit", cb)
		})

		js.Global().Call("makeComponent", "hello-tmp", init)

		hw := goth.CreateElement("hello-tmp")
		goth.AppendChild(cnode, hw)
	}

	// png rendering
	{
		pngb64 := Render()
		content := js.Global().Get("document").Call("getElementById", "content")

		// create div
		img := js.Global().Get("document").Call("createElement", "img")
		img.Set("src", "data:image/png;base64,"+pngb64)

		// append child
		content.Call("appendChild", img)
	}

	forever := make(chan bool)
	<-forever
}

// Render is generate png format by Base64
func Render() string {
	buf := new(bytes.Buffer)

	W := 100
	H := 100
	// var mask *image.Alpha

	dc := gg.NewContext(W, H)
	dc.DrawRectangle(25, 50, 100, 100)
	dc.SetRGBA(255, 128, 0, 1)
	dc.Fill()
	dc.EncodePNG(buf)

	b64 := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, b64)

	io.Copy(encoder, buf)
	encoder.Close()

	return string(b64.Bytes())
}
