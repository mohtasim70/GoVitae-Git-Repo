package main

import "fmt"

type Node struct {
	next *Node
	key  interface{}
}

type List struct {
	head *Node
}

func (L *List) Insert(key interface{}) {
	list := &Node{
		next: L.head,
		key:  key,
	}
	if L.head == nil {
		L.head = list
		return
	}
	list.next = nil

	l := L.head
	for l.next != nil {
		l = l.next
	}
	l.next = list

}
func (L *List) Append(key interface{}) {
	list := &Node{
		next: L.head,
		key:  key,
	}
	if L.head == nil {
		L.head = list
	}

	l := L.head
	for l.next != nil {
		l = l.next
	}
	l.next = L.head
}

func (l *List) Display() {
	list := l.head
	for list != nil {
		fmt.Printf("%+v ->", list.key)
		list = list.next
	}
	fmt.Println()
}

func Display(list *Node) {
	for list != nil {
		fmt.Printf("%v ->", list.key)
		list = list.next
	}
	fmt.Println()
}

func ShowBackwards(list *Node) {
	// for list != nil {
	// 	fmt.Printf("%v <-", list.key)
	// 	list = list.prev
	// }
	// fmt.Println()
}

func (l *List) Reverse() {
	curr := l.head
	var prev *Node
	//	l.tail = l.head

	for curr != nil {
		next := curr.next
		curr.next = prev
		prev = curr
		curr = next
	}
	l.head = prev
	Display(l.head)
}

func main() {
	link := List{}
	// link.Append(5)
	// link.Append(9)
	// link.Append(13)
	// link.Append(22)
	// link.Append(28)
	// link.Append(36)
	link.Insert(5)
	link.Insert(9)
	link.Insert(13)
	link.Insert(22)
	link.Insert(28)
	link.Insert(36)

	fmt.Println("==============================")
	fmt.Printf("Head: %v\n", link.head.key)

	link.Display()
	fmt.Println("==============================")
	fmt.Printf("head: %v\n", link.head.key)

	//link.Reverse()
	fmt.Println("==============================")
}
