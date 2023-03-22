package main

import (
	"strconv"

	"github.com/gregtuc/docker-data-metrics/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/metrics")

	//Fetch all logs
	v1.GET("/", func(c *gin.Context) {
		logs := database.Read()
		if len(logs) >= 1 {
			c.JSON(200, logs)
		} else {
			c.Status(404)
		}
	})

	//Fetch logs filtered by cpu usage under and/or over a specific percent
	v1.GET("/cpu", func(c *gin.Context) {
		var returnLogs []database.Log
		under := c.Query("under")
		over := c.Query("over")
		var underFloat float64
		var overFloat float64

		//Convert query parameters from string to float64
		if under != "" {
			underFloat, err = strconv.ParseFloat(under, 64)
			if err != nil {
				c.Status(400)
				return
			}
		}
		if over != "" {
			overFloat, err = strconv.ParseFloat(over, 64)
			if err != nil {
				c.Status(400)
				return
			}
		}

		//Return logs depending on which parameters were passed
		if under != "" && over != "" {
			logs := database.Read()
			for i := 0; i < len(logs); i++ {
				if logs[i].CPUPercent > overFloat && logs[i].CPUPercent < underFloat {
					returnLogs = append(returnLogs, logs[i])
				}
			}
			c.JSON(200, returnLogs)
		} else if under != "" && over == "" {
			logs := database.Read()
			for i := 0; i < len(logs); i++ {
				if logs[i].CPUPercent < underFloat {
					returnLogs = append(returnLogs, logs[i])
				}
			}
			c.JSON(200, returnLogs)
		} else if under == "" && over != "" {
			logs := database.Read()
			for i := 0; i < len(logs); i++ {
				if logs[i].CPUPercent > overFloat {
					returnLogs = append(returnLogs, logs[i])
				}
			}
			c.JSON(200, returnLogs)
		} else {
			c.Status(400)
		}
	})

	//Fetch logs filtered by timestamp before or after a specific unix time
	v1.GET("/time", func(c *gin.Context) {
		var returnLogs []database.Log
		before := c.Query("before")
		after := c.Query("after")
		var beforeInt int64
		var afterInt int64

		//Convert query parameters from string to int64 (unix)
		if before != "" {
			beforeInt, err = strconv.ParseInt(before, 10, 64)
			if err != nil {
				c.Status(400)
				return
			}
		}
		if after != "" {
			afterInt, err = strconv.ParseInt(after, 10, 64)
			if err != nil {
				c.Status(400)
				return
			}
		}

		//Return logs depending on which timestamps were passed
		if before != "" && after != "" {
			logs := database.Read()
			for i := 0; i < len(logs); i++ {
				if logs[i].Timestamp > afterInt && logs[i].Timestamp < beforeInt {
					returnLogs = append(returnLogs, logs[i])
				}
			}
			c.JSON(200, returnLogs)
		} else if before != "" && after == "" {
			logs := database.Read()
			for i := 0; i < len(logs); i++ {
				if logs[i].Timestamp < beforeInt {
					returnLogs = append(returnLogs, logs[i])
				}
			}
			c.JSON(200, returnLogs)
		} else if before == "" && after != "" {
			logs := database.Read()
			for i := 0; i < len(logs); i++ {
				if logs[i].Timestamp > afterInt {
					returnLogs = append(returnLogs, logs[i])
				}
			}
			c.JSON(200, returnLogs)
		} else {
			c.Status(400)
		}
	})

	//Fetch the latest database log
	v1.GET("/live", func(c *gin.Context) {
		logs := database.Read()
		if len(logs) >= 1 {
			c.JSON(200, logs[len(logs)-1])
		} else {
			c.JSON(404, nil)
		}
	})

	return router
}
