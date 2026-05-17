package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"interface-health-check/database"
	"interface-health-check/models"
)

// 添加巡检接口
func AddAPI(c *gin.Context) {

	var api models.APIInfo

	if err := c.ShouldBindJSON(&api); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&api)

	c.JSON(200, gin.H{"message": "添加成功"})
}

// 删除接口
func DeleteAPI(c *gin.Context) {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	database.DB.Delete(&models.APIInfo{}, id)

	c.JSON(200, gin.H{"message": "删除成功"})
}

// 查看接口列表
func ListAPI(c *gin.Context) {

	var apis []models.APIInfo
	database.DB.Find(&apis)

	c.JSON(200, apis)
}

// 查看巡检结果
func ListChecks(c *gin.Context) {

	var checks []models.APICheck
	database.DB.Order("checked_at desc").Limit(50).Find(&checks)

	c.JSON(200, checks)
}

// Dashboard 页面
func Dashboard(c *gin.Context) {

	var total int64
	var errorCount int64

	database.DB.Model(&models.APICheck{}).Count(&total)
	database.DB.Model(&models.APICheck{}).
		Where("is_error = ?", true).
		Count(&errorCount)

	success := total - errorCount

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"total":   total,
		"success": success,
		"error":   errorCount,
	})
}