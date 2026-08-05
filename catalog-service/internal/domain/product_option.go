package domain

type ProductOptions func(*Product)

func OptionTypes(ots []OptionType) ProductOptions {
	return func(p *Product) {
		p.OptionTypes = ots
	}
}
