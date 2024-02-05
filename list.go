package limit

type node[VALUE any] struct {
	prev  *node[VALUE]
	next  *node[VALUE]
	value VALUE
}

func (n *node[VALUE]) ToTailValues() (result []VALUE) {

	cur := n

	for cur != nil {
		result = append(result, cur.value)
	}

	return
}

type List[VALUE any] struct {
	head *node[VALUE]
	tail *node[VALUE]
}

func (l *List[VALUE]) PutHead(value VALUE) {
	if l.head == nil {
		l.head = &node[VALUE]{
			value: value,
			prev:  nil,
			next:  nil,
		}
		l.tail = l.head
		return
	}

	l.head.prev = &node[VALUE]{
		value: value,
		prev:  nil,
		next:  l.head,
	}

	l.head = l.head.prev
}

func (l *List[VALUE]) IsEmpty() bool {
	return l.head == nil
}

func (l *List[VALUE]) TruncateNodeNext(n *node[VALUE]) *node[VALUE] {
	next := n.next
	n.next = nil
	l.tail = n
	return next
}

// RemoveBack 不做空的判断处理
func (l *List[VALUE]) RemoveBack() {
	if l.tail.prev == nil {
		l.head = nil
		l.tail = nil
		return
	}

	l.tail = l.tail.prev
	l.tail.next = nil
}
