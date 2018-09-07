package goth

import "syscall/js"

// CreateElement is dom ctrl CreateElement
func CreateElement(tagName string) js.Value {
	return js.Global().Get("document").Call("createElement", tagName)
}

// AppendChild is dom ctrl AppendChild
func AppendChild(node js.Value, child js.Value) {
	node.Call("appendChild", child)
}

// RemoveChild is dom ctrl RemoveChild
func RemoveChild(node js.Value, child js.Value) {
	node.Call("removeChild", child)
}

// AddEventListener is dom event RemoveClistenerhild
func AddEventListener(node js.Value, event string, cb js.Callback) {
	node.Call("addEventListener", event, cb)
}

// SetAttribute is dom event SetAttribute
func SetAttribute(node js.Value, name, value string) {
	node.Call("setAttribute", name, value)
}

// RemoveAttribute is dom event RemoveAttribute
func RemoveAttribute(node js.Value, name, value string) {
	node.Call("removeAttribute", name, value)
}
