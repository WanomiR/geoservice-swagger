---
menu:
    after:
        name: binary_tree
        weight: 2
title: Построение сбалансированного бинарного дерева
---

# Задача построить сбалансированное бинарное дерево
Используя AVL дерево, постройте сбалансированное бинарное дерево, на текущей странице.

Нужно написать воркер, который стартует дерево с 5 элементов, и каждые 5 секунд добавляет новый элемент в дерево.

Каждые 5 секунд на странице появляется актуальная версия, сбалансированного дерева.

При вставке нового элемента, в дерево, нужно перестраивать дерево, чтобы оно оставалось сбалансированным.

Как только дерево достигнет 100 элементов, генерируется новое дерево с 5 элементами.

```go
package binary

import (
	"fmt"
	"math/rand"
	"time"
)

type Node struct {
	Key    int
	Height int
	Left   *Node
	Right  *Node
}

type AVLTree struct {
	Root *Node
}

func NewNode(key int) *Node {
	return &Node{Key: key, Height: 1}
}

func (t *AVLTree) Insert(key int) {
	t.Root = insert(t.Root, key)
}

func (t *AVLTree) ToMermaid() string {

}

func height(node *Node) int {

}

func max(a, b int) int {

}

func updateHeight(node *Node) {

}

func getBalance(node *Node) int {

}

func leftRotate(x *Node) *Node {

}

func rightRotate(y *Node) *Node {

}

func insert(node *Node, key int) *Node {

}

func GenerateTree(count int) *AVLTree {

}
```

Не обязательно использовать выше описанный код, можно использовать любую реализацию, выдающую сбалансированное бинарное дерево.

## Mermaid Chart

[MermaidJS](https://mermaid-js.github.io/) is library for generating svg charts and diagrams from text.

## Пример вывода

{{< columns >}}
```tpl
{{</*/* mermaid [class="text-center"]*/*/>}}
graph TD
610 --> 176
176 --> 65
65 --> 48
48 --> 4
65 --> 130
176 --> 340
340 --> 273
273 --> 220
273 --> 317
340 --> 348
348 --> 403
610 --> 758
758 --> 703
703 --> 684
684 --> 673
703 --> 741
741 --> 705
758 --> 949
949 --> 847
{{</*/* /mermaid */*/>}}
```

{{< /columns >}}

{{< mermaid >}}
graph TD
610 --> 176
176 --> 65
65 --> 48
48 --> 4
65 --> 130
176 --> 340
340 --> 273
273 --> 220
273 --> 317
340 --> 348
348 --> 403
610 --> 758
758 --> 703
703 --> 684
684 --> 673
703 --> 741
741 --> 705
758 --> 949
949 --> 847

{{< /mermaid >}}