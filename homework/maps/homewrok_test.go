package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type node struct {
	key   int
	val   int
	left  *node
	right *node
}

func (p *node) insert(newNode *node) {
	if newNode == nil {
		return
	}

	if newNode.key < p.key {
		if p.left != nil {
			p.left.insert(newNode)
			return
		}

		p.left = newNode
		return
	}

	if p.right != nil {
		p.right.insert(newNode)
		return
	}

	p.right = newNode
}

func (p *node) findParent(key int) *node {
	if key < p.key {
		if p.left == nil || p.left.key == key {
			return p
		}
		return p.left.findParent(key)
	}

	if p.right == nil || p.right.key == key {
		return p
	}
	return p.right.findParent(key)
}

func (n *node) forEach(action func(int, int)) {
	if n == nil {
		return
	}

	n.left.forEach(action)
	action(n.key, n.val)
	n.right.forEach(action)
}

type OrderedMap struct {
	root *node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		root: nil,
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	newNode := &node{
		key: key,
		val: value,
	}

	if m.root == nil {
		m.root = newNode
	} else {
		m.root.insert(newNode)
	}

	m.size++
}

func (m *OrderedMap) Erase(key int) {
	if m.root == nil {
		return
	}

	if m.root.key == key {
		root := m.root
		m.root = root.right
		m.root.insert(root.left)
		m.size--
		return
	}

	p := m.root.findParent(key)
	if key < p.key {
		if p.left == nil {
			return
		}

		left := p.left
		p.left = left.right
		p.left.insert(left.left)
		m.size--
		return
	}

	if p.right == nil {
		return
	}

	right := p.right
	p.right = right.left
	p.right.insert(right.right)
	m.size--
}

func (m *OrderedMap) Contains(key int) bool {
	if m.root != nil && m.root.key == key {
		return true
	}

	p := m.root.findParent(key)
	if key < p.key {
		return p.left != nil && p.left.key == key
	}

	return p.right != nil && p.right.key == key
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.root.forEach(action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
