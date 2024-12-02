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
	newItem := &ListItem{Value: v}
	if l.len == 0 {
		listIsZero(l, newItem)
	} else {
		newItem.Next = l.front
		l.front.Prev = newItem
		l.front = newItem
	}
	l.len++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.len == 0 {
		listIsZero(l, newItem)
	} else {
		newItem.Prev = l.back
		l.back.Next = newItem
		l.back = newItem
	}
	l.len++
	return newItem
}

func listIsZero(l *list, newItem *ListItem) {
	l.front = newItem
	l.back = newItem
}

func (l *list) Remove(i *ListItem) {
	if l.len == 0 {
		return
	}
	if i == l.front && i == l.back {
		l.front = nil
		l.back = nil
	} else if i == l.front {
		l.front = i.Next
		l.front.Prev = nil
	} else if i == l.back {
		l.back = i.Prev
		l.back.Next = nil
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.len == 0 || i == l.front {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}
