package data

type Map[K comparable, V any] struct {
	Items []item[K, V]
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{}
}

type item[K comparable, V any] struct {
	key   K
	value V
}

func (m *Map[K, V]) Insert(k K, v V) {
	for i, listItem := range m.Items {
		if listItem.key == k {
			m.Items[i] = item[K, V]{key: k, value: v}
			return
		}
	}
	m.Items = append(m.Items, item[K, V]{
		key:   k,
		value: v,
	})
}

func (m *Map[K, V]) Get(k K) (v V, ok bool) {
	for _, listItem := range m.Items {
		if k == listItem.key {
			v = listItem.value
			ok = true
			break
		}
	}
	return
}

func (m *Map[K, V]) Remove(k K) {
	idx := -1
	for i, listItem := range m.Items {
		if listItem.key == k {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}
	temp := make([]item[K, V], len(m.Items)-1)
	var j int
	for i := range m.Items {
		if i == idx {
			j++
		}
		temp[i] = m.Items[j]
		j++
	}
	m.Items = temp
}
