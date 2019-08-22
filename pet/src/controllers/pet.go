package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//列表数据
func Pets(c *gin.Context) {
	//var Pet model.Pet
	//Pet.Petname = c.Request.FormValue("Petname")
	//Pet.Password = c.Request.FormValue("password")
	//result, err := Pet.Pets()
	//
	//if err != nil {
	c.JSON(http.StatusOK, gin.H{
		"code":    -1,
		"message": "抱歉未找到相关信息",
	})
	return
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code": 1,
	//	"data":   result,
	//})
}

//添加数据
func PetStore(c *gin.Context) {
	//var Pet model.Pet
	//Pet.Petname = c.Request.FormValue("Petname")
	//Pet.Password = c.Request.FormValue("password")
	//id, err := Pet.Insert()

	//if err != nil {
	c.JSON(http.StatusOK, gin.H{
		"code":    -1,
		"message": "添加失败",
	})
	return
	//}
	//c.JSON(http.StatusOK, gin.H{
	//	"code":  1,
	//	"message": "添加成功",
	//	"data":    id,
	//})
}

//修改数据
func PetUpdate(c *gin.Context) {
	//var Pet model.Pet
	//id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	//Pet.Password = c.Request.FormValue("password")
	//result, err := Pet.Update(id)
	//if err != nil || result.ID == 0 {
	c.JSON(http.StatusOK, gin.H{
		"code":    -1,
		"message": "修改失败",
	})
	return
	//}
	//c.JSON(http.StatusOK, gin.H{
	//	"code":  1,
	//	"message": "修改成功",
	//})
}

//删除数据
func PetDestroy(c *gin.Context) {
	//var Pet model.Pet
	//id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	//result, err := Pet.Destroy(id)
	//if err != nil || result.ID == 0 {
	c.JSON(http.StatusOK, gin.H{
		"code":    -1,
		"message": "删除失败",
	})
	return
	//}
	//c.JSON(http.StatusOK, gin.H{
	//	"code":  1,
	//	"message": "删除成功",
	//})
}
