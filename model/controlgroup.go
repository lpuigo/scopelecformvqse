package model

type ControlGroup struct {
	Ref     string
	Label   string
	Section string
	Items   []Item
}
