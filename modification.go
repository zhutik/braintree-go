package braintree

const (
	ModificationKindDiscount = "discount"
	ModificationKindAddOn    = "add_on"
)

type Modification struct {
	Id              string  `xml:"id,omitempty"`
	Amount          float64 `xml:"amount,omitempty"`
	Description     string  `xml:"description,omitempty"`
	Kind            string  `xml:"kind,omitempty"`
	Name            string  `xml:"name,omitempty"`
	NeverExpires    bool    `xml:"never-expires,omitempty"`
	Quantity        int     `xml:"quantity,omitempty"`
	UpdatedAt       string  `xml:"updated_at,omitempty"`
	InheritedFromId string  `xml:"inherited-from-id,omitempty"`
}

type ModificationGroup struct {
	Add *[]Modification `xml:"add>item,omitempty"`
}
