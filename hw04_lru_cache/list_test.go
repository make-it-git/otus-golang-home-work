package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListPush(t *testing.T) {
	t.Run("push front one element", func(t *testing.T) {
		l := NewList()

		item := l.PushFront(123)
		require.Nil(t, item.Next)
		require.Nil(t, item.Prev)
		require.Equal(t, 123, item.Value)
		require.Equal(t, 1, l.Len())
	})

	t.Run("push back one element", func(t *testing.T) {
		l := NewList()

		item := l.PushBack(123)
		require.Nil(t, item.Next)
		require.Nil(t, item.Prev)
		require.Equal(t, 123, item.Value)
		require.Equal(t, 1, l.Len())
	})
}

func TestListRemove(t *testing.T) {
	t.Run("test remove first element", func(t *testing.T) {
		l := NewList()

		item := l.PushFront(123)
		l.PushFront(456)

		require.Equal(t, 2, l.Len())

		front := l.Front()
		back := l.Back()

		require.Equal(t, 456, front.Value)
		require.Nil(t, front.Prev)
		require.Equal(t, back, front.Next)

		require.Equal(t, 123, back.Value)
		require.Nil(t, back.Next)
		require.Equal(t, front, back.Prev)

		l.Remove(item)

		require.Equal(t, 1, l.Len())

		head := l.Front()
		require.Nil(t, head.Next)
		require.Nil(t, head.Prev)
		require.Equal(t, 456, head.Value)
	})

	t.Run("test remove last element", func(t *testing.T) {
		l := NewList()

		l.PushFront(123)
		item := l.PushFront(456)

		require.Equal(t, 2, l.Len())

		front := l.Front()
		back := l.Back()

		require.Equal(t, 456, front.Value)
		require.Nil(t, front.Prev)
		require.Equal(t, back, front.Next)

		require.Equal(t, 123, back.Value)
		require.Nil(t, back.Next)
		require.Equal(t, front, back.Prev)

		l.Remove(item)

		require.Equal(t, 1, l.Len())

		head := l.Front()
		require.Nil(t, head.Next)
		require.Nil(t, head.Prev)
		require.Equal(t, 123, head.Value)
	})

	t.Run("test remove middle element", func(t *testing.T) {
		l := NewList()

		l.PushFront(123)
		middle := l.PushFront(456)
		l.PushFront(789)

		require.Equal(t, 3, l.Len())

		front := l.Front()
		back := l.Back()

		require.Equal(t, 789, front.Value)
		require.Equal(t, middle, front.Next)
		require.Nil(t, front.Prev)

		require.Equal(t, front, middle.Prev)
		require.Equal(t, back, middle.Next)
		require.Equal(t, 456, middle.Value)

		require.Equal(t, 123, back.Value)
		require.Nil(t, back.Next)
		require.Equal(t, middle, back.Prev)

		l.Remove(middle)

		require.Equal(t, 2, l.Len())

		require.Nil(t, middle.Next)
		require.Nil(t, middle.Prev)

		front = l.Front()
		back = l.Back()

		require.Equal(t, 789, front.Value)
		require.Equal(t, 123, back.Value)

		require.Nil(t, back.Next)
		require.Nil(t, front.Prev)

		require.Equal(t, front, back.Prev)
		require.Equal(t, back, front.Next)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{789, 123}, elems)
	})
}

func TestListMove(t *testing.T) {
	t.Run("test move to front first element", func(t *testing.T) {
		l := NewList()

		l.PushFront(123)
		middle := l.PushFront(456)
		first := l.PushFront(789)

		l.MoveToFront(first)

		front := l.Front()
		require.Equal(t, 789, front.Value)
		require.Nil(t, front.Prev)
		require.Equal(t, middle, front.Next)
	})

	t.Run("test move to front middle element", func(t *testing.T) {
		l := NewList()

		l.PushFront(1)
		middle := l.PushFront(2)
		first := l.PushFront(3)

		l.MoveToFront(middle)

		front := l.Front()
		back := l.Back()

		require.Equal(t, 2, front.Value)
		require.Nil(t, front.Prev)
		require.Equal(t, first, front.Next)

		require.Equal(t, 1, back.Value)
		require.Nil(t, back.Next)
		require.Equal(t, first, back.Prev)
	})

	t.Run("test move to front last element", func(t *testing.T) {
		l := NewList()

		last := l.PushFront(1)
		l.PushFront(2)
		first := l.PushFront(3)

		// 3 2 1
		l.MoveToFront(last)
		// 1 3 2

		front := l.Front()
		back := l.Back()

		require.Equal(t, 1, front.Value)
		require.Nil(t, front.Prev)
		require.Equal(t, first, front.Next)

		require.Equal(t, 2, back.Value)
		require.Nil(t, back.Next)
		require.Equal(t, first, back.Prev)
	})
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		x1 := l.PushFront(10) // [10]
		x2 := l.PushBack(20)  // [10, 20]
		x3 := l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		require.Nil(t, x1.Prev)
		require.Equal(t, x1.Next, x2)
		require.Equal(t, x2.Prev, x1)
		require.Equal(t, x2.Next, x3)
		require.Nil(t, x3.Next)
		require.Equal(t, x3.Prev, x2)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{10, 20, 30}, elems)

		middle := l.Front().Next // 20
		require.Equal(t, 20, middle.Value)
		l.Remove(middle) // [10, 30]
		require.Nil(t, middle.Next)
		require.Nil(t, middle.Prev)
		require.Equal(t, 2, l.Len())
		front := l.Front()
		back := l.Back()
		require.Equal(t, 10, front.Value)
		require.Nil(t, front.Prev)
		require.Equal(t, back, front.Next)
		require.Equal(t, 30, back.Value)
		require.Equal(t, front, back.Prev)
		require.Nil(t, back.Next)

		elems = make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{10, 30}, elems)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		elems = make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{80, 60, 40, 10, 30, 50, 70}, elems)

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems = make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
