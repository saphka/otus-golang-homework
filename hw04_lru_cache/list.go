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
	len        int
	head, tail *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newHead := &ListItem{
		Value: v,
	}
	return l.pushFrontInternal(newHead)
}

func (l *list) pushFrontInternal(newHead *ListItem) *ListItem {
	oldHead := l.head
	l.head = newHead
	l.head.Next = oldHead
	l.head.Prev = nil
	if oldHead != nil {
		oldHead.Prev = l.head
	} else {
		l.tail = l.head
	}
	l.len++
	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	oldTail := l.tail
	l.tail = &ListItem{
		Value: v,
		Prev:  oldTail,
	}
	if oldTail != nil {
		oldTail.Next = l.tail
	} else {
		l.head = l.tail
	}
	l.len++
	return l.tail
}

func (l *list) Remove(i *ListItem) {
	switch {
	case i.Prev != nil && i.Next != nil: // middle element
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	case i.Prev != nil: // tail
		l.tail = i.Prev
		l.tail.Next = nil
	case i.Next != nil: // head
		l.head = i.Next
		l.head.Prev = nil
	default: // single element
		l.head = nil
		l.tail = nil
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.pushFrontInternal(i)
}

func NewList() List {
	return new(list)
}
