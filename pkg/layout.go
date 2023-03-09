package pkg

import "github.com/google/uuid"

const ID = "ID"

type Layout struct {
	positions map[string]position
}

type position struct {
	X, Y int
}

func NewLayout() Layout {
	return Layout{make(map[string]position)}
}

func (l *Layout) horizontalLayout(root Node, depth int, id string) {
	if root == nil {
		return
	}

	root[ID] = id

	if _, ok := l.positions[id]; !ok {
		l.positions[id] = position{depth, 0}
		root[X_POSITION_FACTOR] = float64(depth)
		root[Y_POSITION_FACTOR] = float64(0)
	}
	if root[PLANS_PROP] == nil {
		return
	}

	children := root[PLANS_PROP].([]interface{})
	// calculate positions of children recursively
	x := l.positions[id].X - 1               // move left
	y := l.positions[id].Y - len(children)/2 // center vertically
	for _, child := range children {
		node := child.(Node)
		childId := uuid.New().String()
		node[ID] = childId
		l.positions[childId] = position{x, y}
		node[X_POSITION_FACTOR] = float64(x)
		node[Y_POSITION_FACTOR] = float64(y)
		y++
		l.horizontalLayout(node, depth+1, childId)
	}
}

func (l *Layout) PrintTreeHorizontal(root Node) {
	id := uuid.New().String()
	l.horizontalLayout(root, 0, id)
}
