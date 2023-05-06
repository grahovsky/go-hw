package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.front = &ListItem{Value: v, Next: l.front, Prev: nil}
	if l.len == 0 {
		l.back = l.front
	} else {
		l.front.Next.Prev = l.front
	}

	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.back = &ListItem{Value: v, Next: nil, Prev: l.back}
	if l.len == 0 {
		l.front = l.back
	} else {
		l.back.Prev.Next = l.back
	}

	l.len++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.front = l.front.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.back = l.back.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	i.Next = nil
	i.Prev = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	switch {
	case i.Prev == nil:
		return
	case i.Next == nil:
		l.back = i.Prev
	default:
		i.Next.Prev = i.Prev
	}

	l.front.Prev, l.front, i.Next, i.Prev.Next = i, i, l.front, i.Next
}

func NewList() List {
	return new(list)
}
