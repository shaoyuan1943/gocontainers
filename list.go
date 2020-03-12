package gocontainers

type Entity struct {
	prev *Entity     // prev
	next *Entity     // next
	l    *List       // list
	v    interface{} // value
}

func (e *Entity) Init(value interface{}) {
	e.v = value
}

func (e *Entity) Value() interface{} {
	return e.v
}

func (e *Entity) Next() *Entity {
	if p := e.next; e.l != nil && p != e.l.root {
		return p
	}
	return nil
}

func (e *Entity) Prev() *Entity {
	if p := e.prev; e.l != nil && p != e.l.root {
		return p
	}

	return nil
}

type List struct {
	root *Entity
	len  int
}

func NewList() *List {
	l := &List{
		root: &Entity{},
		len:  0,
	}

	l.root.next = l.root
	l.root.prev = l.root
	return l
}

func (l *List) Len() int {
	return l.len
}

func (l *List) Front() *Entity {
	if l.len == 0 {
		return nil
	}

	return l.root.next
}

func (l *List) Back() *Entity {
	if l.len == 0 {
		return nil
	}

	return l.root.prev
}

func (l *List) insert(e *Entity, at *Entity, force bool) *Entity {
	if e.l != nil && !force {
		return nil
	}

	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.l = l

	l.len++
	return e
}

func (l *List) remove(e *Entity) *Entity {
	e.prev.next = e.next
	e.next.prev = e.prev

	e.next = nil
	e.prev = nil
	e.l = nil

	l.len--
	return e
}

func (l *List) move(e *Entity, at *Entity) *Entity {
	if e == at {
		return e
	}

	e.prev.next = e.next
	e.next.prev = e.prev

	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e

	return e
}

func (l *List) Remove(e *Entity) *Entity {
	if e == nil || e.l != l {
		return nil
	}

	return l.remove(e)
}

func (l *List) PushFront(e *Entity) *Entity {
	return l.insert(e, l.root, false)
}

func (l *List) PushBack(e *Entity) *Entity {
	return l.insert(e, l.root.prev, false)
}

func (l *List) InsertBefore(e *Entity, mark *Entity) *Entity {
	if mark.l != l {
		return nil
	}

	return l.insert(e, mark.prev, false)
}

func (l *List) InsertAfter(e *Entity, mark *Entity) *Entity {
	if mark.l != l {
		return nil
	}

	return l.insert(e, mark.next, false)
}

func (l *List) MoveToFront(e *Entity) {
	if e.l != l || l.root.next == e {
		return
	}

	l.move(e, l.root)
}

func (l *List) MoveToBack(e *Entity) {
	if e.l != l || l.root.prev == e {
		return
	}

	l.move(e, l.root.prev)
}

func (l *List) MoveBefore(e, mark *Entity) {
	if e.l != l || e == mark || mark.l != l {
		return
	}

	l.move(e, mark.prev)
}

func (l *List) MoveAfter(e, mark *Entity) {
	if e.l != l || e == mark || mark.l != l {
		return
	}

	l.move(e, mark)
}

func (l *List) PushBackList(other *List) {
	e1 := other.Front()
	for {
		if e1 == nil {
			break
		}

		e2 := *e1
		l.insert(&e2, l.root.prev, true)
		e1 = e1.Next()
	}
}

func (l *List) PushFrontList(other *List) {
	e1 := other.Back()
	for {
		if e1 == nil {
			break
		}

		e2 := *e1
		l.insert(&e2, l.root, true)
		e1 = e1.Prev()
	}
}

func (l *List) PopFront() *Entity {
	e := l.Front()
	if e != nil {
		l.Remove(e)
	}
	return e
}

func (l *List) PopBack() *Entity {
	e := l.Back()
	if e != nil {
		l.Remove(e)
	}
	return e
}

func (l *List) Clear() {
	l.root.next = l.root
	l.root.prev = l.root
	l.len = 0
}
