package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"syscall/js"
	"time"

	"github.com/fogleman/gg"
)

func main() {
	fmt.Println("Hello wasm")

	{
		var cb js.Callback

		cb = js.NewCallback(func(args []js.Value) {
			fmt.Println("button clicked")

			// callbackを閉じる。 once
			cb.Release()
		})
		js.Global().Get("document").Call("getElementById", "myButton").Call("addEventListener", "click", cb)
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
			content := js.Global().Get("document").Call("getElementById", "content")

			// create div
			div := js.Global().Get("document").Call("createElement", "div")
			div.Set("textContent", "new row")

			// append child
			content.Call("appendChild", div)

			// Releaseしていないのでclickのたびに実行する
			// cb.Release()
		})
		js.Global().Get("document").Call("getElementById", "createDiv").Call("addEventListener", "click", cb)
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
	// for {
	// 	time.Sleep(time.Second)
	// 	fmt.Println(fmt.Sprintf("Hello wasm outer loop At: %d", counter))
	// 	if counter > 10 {
	// 		fmt.Println("break outer loop")
	// 		break
	// 	}
	// 	counter++
	// }

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

// base 64 pngの生成
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
