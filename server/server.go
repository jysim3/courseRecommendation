package main

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/labstack/echo"
)

// {
// 	"SUBJECT_AREA":{
// 	  "COURS_CODE1": {
// 		"available_sems": [1,2,3],
// 		"prereqs": ["COURSE_CODEX"],
// 		"excluded": ["COURSE_CODEX", "COURSE_CODEY"],
// 		"UOC": 6,
// 		"GENED": true
// 	  },
// 	  "COURSE_CODE2": {
// 		"available_sems": [1],
// 		"prereqs": ["COURSE_CODEX", "COURSE_CODEY"],
// 		"excluded": ["COURSE_CODEX", "COURSE_CODEY"],
// 		"UOC": 12,
// 		"GENED": false
// 	  }
// 	}
//   }

func main() {
	e := echo.New()


	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	logrus.SetFormatter(logrus.Form)

	e.POST("/query", h echo.HandlerFunc, m ...echo.MiddlewareFunc)

	e.Logger.Fatal(e.Start(":1323"))
}
