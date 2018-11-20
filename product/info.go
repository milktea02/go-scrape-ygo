package product

type Info struct {
	Name     string
	Set      string
	Variants []*Variant
}

type Variant struct {
	Condition string
	Price     float64
	Quantity  int64
}
