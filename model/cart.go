package model

type (
	Product struct {
		Id          int
		Name        string
		Description string
		Image       string
		Price       float32
	}
)

type (
	CartItem struct {
		ProductId int
		Quantity  int
	}

	CartState struct {
		Items     []CartItem `json:",omitempty"`
		TrxID     string     `json:",omitempty"`
		PaymentID string     `json:",omitempty"`
		Checkout  bool       `json:",omitempty"`
		Payment   bool       `json:",omitempty"`
	}
)

type AddToCartSignal struct {
	Route string
	Item  CartItem
}

type RemoveFromCartSignal struct {
	Route string
	Item  CartItem
}

type CheckoutSignal struct {
	Route string
}

type PaymentSignal struct {
	Route string
}

type Transaction struct {
	TrxID     string `json:",omitempty"`
	PaymentID string `json:",omitempty"`
}
