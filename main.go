package main

import (
	"fmt"
	"syscall/js"
	"time"
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
	for {
		time.Sleep(time.Second)
		fmt.Println(fmt.Sprintf("Hello wasm outer loop At: %d", counter))
		if counter > 10 {
			fmt.Println("break outer loop")
			break
		}
		counter++
	}

	forever := make(chan bool)
	<-forever
}
