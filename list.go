package limit

type node[VALUE any] struct {
	prev  *node[VALUE]
	next  *node[VALUE]
	value VALUE
}

type List[VALUE any] struct {
	size int
	head *node[VALUE]
	tail *node[VALUE]

	zero VALUE
}

func (l *List[VALUE]) PutHead(value VALUE) {
	l.size++
	n := &node[VALUE]{
		value: value,
		prev:  nil,
		next:  nil,
	}

	if l.head == nil {
		l.head = n
		l.tail = l.head
		return
	}

	l.head.prev = n
	n.next = l.head
	l.head = n
}

func (l *List[VALUE]) IsEmpty() bool {
	return l.head == nil
}

// RemoveBack 不做空的判断处理
func (l *List[VALUE]) RemoveBack() VALUE {

	if l.size == 0 {
		return l.zero
	}

	l.size--

	if l.tail.prev == nil {
		n := l.tail
		l.head = nil
		l.tail = nil
		return n.value
	}
	n := l.tail

	l.tail = l.tail.prev
	l.tail.next = nil

	return n.value
}
