package goth

import "syscall/js"

// CreateElement is dom ctrl CreateElement
func CreateElement(tagName string) Node {
	return Node{
		jsv: js.Global().Get("document").Call("createElement", tagName),
	}
}

// AppendChildNode is dom ctrl AppendChildNode
func AppendChild(node Node, child Node) {
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

func (n Node) Get(prop string) js.Value {
	return n.jsv.Get(prop)
}

func (n Node) SetAttribute(name, value string) Node {
	SetAttribute(n.jsv, name, value)
	return n
}

func (n Node) AddEventListener(event string, cb js.Callback) Node {
	AddEventListener(n.jsv, event, cb)
	return n
}

func (n Node) AttachShadow() Node {
	obj := make(map[string]interface{})
	obj["mode"] = "open"
	n.jsv.Call("attachShadow", js.ValueOf(obj))
	return n
}

func (n Node) CreateShadowRoot() Node {
	return Node{jsv: n.jsv.Call("createShadowRoot")}
}

func (n Node) AppendChild(child Node) Node {
	AppendChild(n, child)
	return n
}
