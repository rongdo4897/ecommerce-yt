package database

import "errors"

var (
	ErrCantFindProduct     = errors.New("can't find the product")
	ErrCantDecoderProducts = errors.New("can't find decoder product")
	ErrUserIdIsNotValid    = errors.New("user is not valid")
	ErrCantUpdateUser      = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart  = errors.New("can't remove item from the cart")
	ErrCantGetItem         = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem     = errors.New("can't update the purchase")
)

func AddProductToCart() {

}

func RemoveCarytItem() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
