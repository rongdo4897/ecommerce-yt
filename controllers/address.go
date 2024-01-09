package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rongdo4897/ecommerce-yt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {

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
