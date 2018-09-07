package goth

import "syscall/js"

// CreateElement is dom ctrl CreateElement
func CreateElement(tagName string) js.Value {
	return js.Global().Get("document").Call("createElement", tagName)
}

// CreateElement is dom ctrl CreateElement
func CreateElementNode(tagName string) Node {
	return Node{
		jsv: js.Global().Get("document").Call("createElement", tagName),
	}
}

// AppendChild is dom ctrl AppendChild
func AppendChild(node js.Value, child js.Value) {
	node.Call("appendChild", child)
}

// AppendChildNode is dom ctrl AppendChildNode
func AppendChildNode(node Node, child Node) {
	node.jsv.Call("appendChild", child.jsv)
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

func NewNode(jsv js.Value) Node {
	return Node{jsv: jsv}
}

type Node struct {
	jsv js.Value
}

func (n Node) Set(prop, value string) Node {
	n.jsv.Set(prop, value)
	return n
}

func (n Node) SetAttribute(name, value string) Node {
	SetAttribute(n.jsv, name, value)
	return n
}

func (n Node) AddEventListener(event string, cb js.Callback) Node {
	AddEventListener(n.jsv, event, cb)
	return n
}
