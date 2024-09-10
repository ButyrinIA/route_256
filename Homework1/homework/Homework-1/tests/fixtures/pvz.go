package fixtures

import (
	"homework/Homework-1/internal/pkg/repository"
	"homework/Homework-1/tests/states"
)

type PVZBuilder struct {
	instance *repository.PVZ
}

func PVZ() *PVZBuilder {
	return &PVZBuilder{instance: &repository.PVZ{}}
}

func (b *PVZBuilder) ID(v int64) *PVZBuilder {
	b.instance.ID = v
	return b
}

func (b *PVZBuilder) Name(v string) *PVZBuilder {
	b.instance.Name = v
	return b
}

func (b *PVZBuilder) Address(v string) *PVZBuilder {
	b.instance.Address = v
	return b
}

func (b *PVZBuilder) Contact(v string) *PVZBuilder {
	b.instance.Contact = v
	return b
}

func (b *PVZBuilder) P() *repository.PVZ {
	return b.instance
}

func (b *PVZBuilder) V() repository.PVZ {
	return *b.instance
}

func (b *PVZBuilder) Valid() *PVZBuilder {
	return PVZ().ID(states.PVZ1ID).Name(states.PVZ1Name).Address(states.PVZ1Address).Contact(states.PVZ1Contact)
}
