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

type Layout struct {
}

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
		button := goth.CreateElementNode("button").
			SetAttribute("id", "myButton").
			Set("innerHTML", "Go Callback Test").
			AddEventListener("click", cb)
		goth.AppendChildNode(cnode, button)
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
			goth.AppendChild(content, div)

			// Releaseしていないのでclickのたびに実行する
			// cb.Release()
		})
		button := goth.CreateElement("button")
		goth.SetAttribute(button, "id", "createDiv")
		button.Set("innerHTML", "Go Create DOM")
		goth.AddEventListener(button, "click", cb)
		goth.AppendChild(content, button)
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
			element := i[0]
			element.Set("innerHTML", "Hello world!")
			fmt.Println("format ppp", element)
		})

		js.Global().Call("makeComponent", "hello-world", init)

		hw := goth.CreateElementNode("hello-world")
		goth.AppendChildNode(cnode, hw)
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

// Generate png format by Base64
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
