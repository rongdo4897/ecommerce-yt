package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
func main() {
	// Tạo client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Tạo context với timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tạo kết nối đến MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Kiểm tra kết nối
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Đóng kết nối khi không cần sử dụng nữa
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}
*/

var Client *mongo.Client = DBSet()

func DBSet() *mongo.Client {
	// Tạo 1 client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	// tạo một context với một thời gian chờ (timeout) 10 giây.
	// Bất cứ hoạt động nào được thực hiện trong ngữ cảnh này sẽ bị hủy bỏ nếu nó không hoàn thành trong khoảng thời gian này.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tạo kết nối
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Kiểm tra kết nối tới mongo
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("failed to connect to mongodb: ", err)
		return nil
	}

	fmt.Println("Successfully connected to mongodb")

	return client
}

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	// truy cập một bảng (collection) cụ thể trong một database cụ thể.
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return productCollection
}
