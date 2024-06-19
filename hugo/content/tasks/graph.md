---
menu:
    after:
        name: graph
        weight: 1
title: Построение графа
---

# Построение графа

Нужно написать воркер, который будет строить граф на текущей странице, каждые 5 секунд
От 5 до 30 элементов, случайным образом. Все ноды графа должны быть связаны.

```go
type Node struct {
    ID int
    Name string
	Form string // "circle", "rect", "square", "ellipse", "round-rect", "rhombus"
    Links []*Node
}
```

## Mermaid Chart

[MermaidJS](https://mermaid-js.github.io/) is library for generating svg charts and diagrams from text.

## Пример

{{< columns >}}
```tpl
{{</*/* mermaid [class="text-center"]*/*/>}}
graph LR
A[Square Rect] --> B((Circle))
A --> C(Round Rect)
B --> D{Rhombus}
C --> D
C --> B
{{</*/* /mermaid */*/>}}
```

<--->

{{< mermaid >}}
graph LR
A[Square Rect] --> B((Circle))
A --> C(Round Rect)
B --> D{Rhombus}
C --> D
C --> B
{{< /mermaid >}}

{{< /columns >}}