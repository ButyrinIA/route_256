package storage

type PackagingType struct {
	Name         string
	MaxWeight    int
	CostIncrease int
}

// Карта для хранения информации о каждом виде упаковки
var packagingType = map[string]PackagingType{
	"Packet": {Name: "Packet", MaxWeight: 10, CostIncrease: 5},
	"Box":    {Name: "Box", MaxWeight: 30, CostIncrease: 20},
	"Film":   {Name: "Film", MaxWeight: 0, CostIncrease: 1},
}
