package eliminate

type Node struct {
	Key  string
	pre  *Node
	next *Node
}

type LRUCache struct {
	limit   int
	HashMap map[string]*Node
	head    *Node
	end     *Node
}

func Constructor(capacity int) LRUCache {
	lruCache := LRUCache{limit: capacity}
	lruCache.HashMap = make(map[string]*Node, capacity)
	return lruCache
}

func (l *LRUCache) Get(key string) {
	if v, ok := l.HashMap[key]; ok {
		l.refreshNode(v)
	}
}

func (l *LRUCache) Put(key string) {
	if v, ok := l.HashMap[key]; ok {
		l.refreshNode(v)
		return
	}

	if len(l.HashMap) >= l.limit {
		oldKey := l.removeNode(l.head)
		delete(l.HashMap, oldKey)
	}
	node := Node{Key: key}
	l.addNode(&node)
	l.HashMap[key] = &node

}

func (l *LRUCache) refreshNode(node *Node) {
	if node == l.end {
		return
	}
	l.removeNode(node)
	l.addNode(node)
}

func (l *LRUCache) removeNode(node *Node) string {
	if node == l.end {
		l.end = l.end.pre
		l.end.next = nil
	} else if node == l.head {
		l.head = l.head.next
		l.head.pre = nil
	} else {
		node.pre.next = node.next
		node.next.pre = node.pre
	}
	return node.Key
}

func (l *LRUCache) addNode(node *Node) {
	if l.end != nil {
		l.end.next = node
		node.pre = l.end
		node.next = nil
	}
	l.end = node
	if l.head == nil {
		l.head = node
	}
}

/**
 * 内存淘汰
 */
func (l *LRUCache) Remove() {
	//todo

}

func (l *LRUCache) Method() {

}
