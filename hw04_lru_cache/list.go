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
	head   *ListItem
	tail   *ListItem
	length int
}

func (l list) Len() int {
	return l.length
}

func (l list) Front() *ListItem {
	return l.head
}

func (l list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	l.length++

	if l.head == nil {
		l.head = item
		l.tail = item
		return l.head
	}

	//nolint:ifshort
	oldHead := l.head
	item.Next = l.head
	l.head = item
	if oldHead != nil {
		oldHead.Prev = l.head
	}
	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	l.length++

	//nolint:ifshort
	oldTail := l.tail
	if l.tail != nil {
		l.tail.Next = item
	}

	l.tail = item
	if oldTail != nil {
		oldTail.Next = item
		item.Prev = oldTail
	}

	return l.tail
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	l.length--

	if l.head == i {
		l.head = l.head.Next
	}

	if l.tail == i {
		l.tail = l.tail.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Prev = nil
	i.Next = nil
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	if l.head == i {
		return
	}

	if l.tail == i {
		l.tail = i.Prev
	}

	oldHead := l.head
	l.head = i
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	i.Next = oldHead
	i.Prev = nil
	oldHead.Prev = l.head
}

func NewList() List {
	return new(list)
}
