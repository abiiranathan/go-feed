package routes

import (
	"log"
	"strconv"

	"github.com/abiiranathan/gin_api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginatedFeeds is a function that returns the paginated feeds
func PaginatedFeeds(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// extract the query param
		query := c.Query("query")

		// get the page number as integer
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}

		// get the limit as integer
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}

		// fetch all the feeds and paginate them with our page and limit
		var feeds []models.Feed

		pagination := models.Paginate(db, &feeds, page, limit, query)

		// send the paginated feeds
		c.JSON(200, gin.H{
			"data":       feeds,
			"pagination": pagination,
		})

	}
}

// GET all feeds route
func GetFeeds(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var feeds []models.Feed
		if err := db.Find(&feeds).Error; err != nil {
			c.AbortWithStatus(500)
			log.Println(err)
		} else {
			c.JSON(200, feeds)
		}
	}
}

// GET a feed route
func GetFeed(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var feed models.Feed
		if err := db.First(&feed, id).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		} else {
			c.JSON(200, feed)
		}
	}
}

// POST a feed route
func PostFeed(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var feed models.Feed
		if err := c.BindJSON(&feed); err != nil {
			c.AbortWithStatus(400)
			log.Println(err)
		} else {
			if err := db.Create(&feed).Error; err != nil {
				// send json with error message
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
				log.Println(err)
			} else {
				c.JSON(201, feed)
			}
		}
	}
}

// PUT a feed route
func PutFeed(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var feed models.Feed
		id := c.Param("id")
		if err := db.First(&feed, id).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		} else {
			if err := c.BindJSON(&feed); err != nil {
				// send json with error message
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
				log.Println(err)
			} else {
				db.Save(&feed)
				c.JSON(200, feed)
			}
		}
	}
}

// DELETE a feed route
func DeleteFeed(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var feed models.Feed
		if err := db.First(&feed, id).Error; err != nil {
			c.AbortWithStatus(404)
			log.Println(err)
		} else {
			db.Delete(&feed)
			c.JSON(200, gin.H{
				"message": "feed deleted",
			})
		}
	}
}

// match feeds matching the search query of title and description
func SearchFeeds(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var feeds []models.Feed
		query := c.Query("query")

		if err := db.Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&feeds).Error; err != nil {
			c.AbortWithStatus(500)
			log.Println(err)
		} else {
			c.JSON(200, feeds)
		}
	}
}

// setupRoutes is a function that sets up the routes for the API
// Grouping the routes together
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	feedRoutes := r.Group("/feeds")

	feedRoutes.GET("/", PaginatedFeeds(db))
	feedRoutes.GET("/:id", GetFeed(db))
	feedRoutes.POST("/", PostFeed(db))
	feedRoutes.PUT("/:id", PutFeed(db))
	feedRoutes.DELETE("/:id", DeleteFeed(db))
}
