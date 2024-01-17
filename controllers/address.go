package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rongdo4897/ecommerce-yt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var addressModel models.Address
		addressModel.Address_ID = primitive.NewObjectID()
		if err = c.BindJSON(&addressModel); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		/*
			$match: Đây là một giai đoạn để lọc các tài liệu theo một điều kiện nhất định, Trong trường hợp này, tài liệu sẽ được lọc theo trường "_id" bằng giá trị "address".
			Điều này tạo một đối tượng match_filter kiểu bson.D (bson.Document) với một phần tử có key là "$match" và giá trị là một bson.D khác chứa điều kiện lọc.
			Điều này tạo ra một điều kiện $match trong pipeline để lọc các tài liệu dựa trên giá trị "_id" bằng address.
		*/
		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		/*
			$unwind: Giai đoạn này được sử dụng để chia các giá trị trong mảng thành các tài liệu độc lập.
			Trong trường hợp này, có vẻ như "address" là một mảng, và $unwind sẽ "mở rộng" nó, tạo ra các bản ghi riêng lẻ cho mỗi phần tử trong mảng.
			Trong đoạn mã này, nó mở rộng mảng có tên "address" thành các bản ghi độc lập. path là key chứa mảng mà bạn muốn mở rộng, và "$address" là đường dẫn của mảng đó trong tài liệu
		*/
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		/*
			$group: Đây là giai đoạn nhóm, nơi bạn có thể thực hiện các phép toán nhóm như tổng, đếm, lấy giá trị lớn nhất, lấy giá trị nhỏ nhất, v.v.
			Trong trường hợp này, tài liệu được nhóm dựa trên giá trị "_id" là "$address_id", và sau đó, cho mỗi nhóm, đếm số lượng bằng cách sử dụng $sum.
		*/
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		/*
		 UserCollection.Aggregate thực hiện toàn bộ pipeline trên collection UserCollection sử dụng các giai đoạn đã được xây dựng
		 , và kết quả được trả về dưới dạng một con trỏ cursor pointCursor.
		*/
		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		/*
			addressInfo là một slice (mảng động) chứa các tài liệu kết quả của truy vấn MongoDB.
			bson.M là một kiểu dữ liệu trong thư viện BSON của MongoDB, đại diện cho một tài liệu BSON dưới dạng map[string]interface{}.
			Trong Go, bson.M thường được sử dụng khi không biết chính xác cấu trúc của tài liệu hoặc muốn làm việc với dữ liệu không cố định.
		*/
		var addressInfo []bson.M
		/*
			pointCursor là con trỏ cursor chứa kết quả từ truy vấn aggregation.
			Phương thức All của cursor được sử dụng để lấy tất cả các tài liệu từ cursor và đổ vào slice addressInfo.
			ctx là context được sử dụng trong truy vấn, và &addressInfo là địa chỉ của slice để lưu trữ kết quả.
		*/
		if err = pointCursor.All(ctx, &addressInfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addressModel}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}

			//TODO: Thiếu xử lý json success
			// c.IndentedJSON(200, "Add address successfully")
		} else {
			c.IndentedJSON(400, "Not Allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return
		}

		listAddress := make([]models.Address, 0)
		/*
			tạo một đối tượng primitive.ObjectID từ một chuỗi hex string user_id
			Trong MongoDB, ObjectID là một kiểu dữ liệu đặc biệt thường được sử dụng để biểu diễn trường "_id" của các tài liệu.

			primitive.ObjectID: Là kiểu dữ liệu trong thư viện MongoDB Go (mgo) để biểu diễn một ObjectID, là một loại dữ liệu duy nhất được sử dụng cho trường "_id" trong MongoDB.
			ObjectIDFromHex(user_id): Là một hàm tạo (constructor) của thư viện MongoDB Go, nhận một chuỗi hex user_id làm đối số và trả về một đối tượng ObjectID. Hàm này chuyển đổi chuỗi hex thành một đối tượng ObjectID.
		*/
		userObjectId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		/*
			tạo một bộ lọc MongoDB dựa trên giá trị của trường "_id" (ID) trong một tài liệu

			bson.D: Là kiểu dữ liệu trong thư viện MongoDB Go (mgo) để biểu diễn một tài liệu BSON (Binary JSON) dưới dạng slice của các phần tử primitive.E (Element).
			primitive.E: Là một cặp khóa-giá trị, đại diện cho một phần tử trong một tài liệu BSON. Trong trường hợp này, Key là khóa và Value là giá trị.
			{Key: "_id", Value: userObjectId}: Đây là một cặp khóa-giá trị trong một tài liệu BSON.
			Key: "_id" chỉ định trường "_id" của tài liệu MongoDB.
			Value: userObjectId là giá trị mà bạn muốn sử dụng để so khớp với trường "_id" trong tài liệu MongoDB. Giả sử userObjectId là một giá trị ID mong muốn.

		*/
		filter := bson.D{primitive.E{Key: "_id", Value: userObjectId}}
		/*
			tạo một tài liệu BSON (Binary JSON) chứa một phần tử $set trong ngữ cảnh cập nhật MongoDB.

			bson.D: Đây là kiểu dữ liệu trong thư viện MongoDB Go (mgo) để biểu diễn một tài liệu BSON dưới dạng slice của các phần tử primitive.E (Element).
			primitive.E: Là một cặp khóa-giá trị, đại diện cho một phần tử trong một tài liệu BSON. Trong trường hợp này, Key là khóa và Value là giá trị.

			Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: listAddress}}: Đây là một cặp khóa-giá trị trong một tài liệu BSON.
			Key: "$set" chỉ định một toán tử $set, thường được sử dụng trong truy vấn cập nhật MongoDB để thiết lập giá trị mới cho một trường cụ thể.
			Value: bson.D{primitive.E{Key: "address", Value: listAddress}} chỉ định rằng trường "address" sẽ được thiết lập bằng giá trị của listAddress. Đây có thể là một giá trị cụ thể hoặc một tài liệu BSON khác.
		*/
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: listAddress}}}}
		/*
			thực hiện một truy vấn cập nhật (update) trên một bộ sưu tập (collection) MongoDB.

			UserCollection: Đây có thể là biến hoặc đối tượng đại diện cho bộ sưu tập MongoDB. Trong ngữ cảnh này, đây có thể là biến mô tả bộ sưu tập "users" trong cơ sở dữ liệu.
			UpdateOne: Là một phương thức của đối tượng bộ sưu tập MongoDB (UserCollection) để thực hiện một truy vấn cập nhật trên một tài liệu.
			ctx: Là đối tượng context.Context, thường được sử dụng để quản lý các yêu cầu hủy bỏ, giới hạn thời gian chờ, và các giá trị context khác.
			filter: Là một đối tượng BSON (Binary JSON) hoặc một biểu thức có thể đánh giá thành đúng hoặc sai, xác định điều kiện để xác định tài liệu cần cập nhật.
			update: Là một đối tượng BSON chứa các toán tử cập nhật và giá trị mới cho các trường cần được cập nhật trong tài liệu.
		*/
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Wrong command")
			return
		}
		defer cancel()
		ctx.Done()

		c.IndentedJSON(200, "Successfully deleted")
	}
}
