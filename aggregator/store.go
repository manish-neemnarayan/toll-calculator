package main

import "github.com/manish-neemnarayan/toll-calculator/types"

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}

func (m *MemoryStore) Get(obuId int) (inv float64, err error) {
	inv, ok := m.data[obuId]
	if !ok {
		inv = 0.0
		return inv, err
	}
	return inv, nil
}
