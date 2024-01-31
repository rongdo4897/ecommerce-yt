package database

import (
	"context"
	"errors"
	"log"

	"github.com/rongdo4897/ecommerce-yt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct     = errors.New("can't find the product")
	ErrCantDecoderProducts = errors.New("can't find decoder product")
	ErrUserIdIsNotValid    = errors.New("user is not valid")
	ErrCantUpdateUser      = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart  = errors.New("can't remove item from the cart")
	ErrCantGetItem         = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem     = errors.New("can't update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDB, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser
	err = searchFromDB.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecoderProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
