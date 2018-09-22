package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
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

var resultsTmpl = template.Must(template.New("result.tmpl").ParseFiles("./html_file/result.tmpl"))

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{})

	f, err := os.Open("./courseinfo.json")
	if err != nil {
		logrus.Fatal(err)
	}

	var courses AllCourses
	err = json.NewDecoder(f).Decode(&(courses.Courses))
	if err != nil {
		f.Close()
		logrus.Fatal(err)
	}

	f.Close()

	e := echo.New()

	e.POST("/", func(c echo.Context) error {
		var user User

		for i := 1; i < 20; i++ {
			val := c.FormValue("course" + strconv.Itoa(i))
			if val == "" {
				break
			}

			user.Completed = append(user.Completed, strings.Split(val, " ")[0])
		}

		faculty := strings.Split(c.FormValue("search_input"), ":")[0]
		facultyFilter := func(c *Course) bool {
			return strings.HasPrefix(c.ID, faculty)
		}

		suggestions := courses.Suggestions(&user, facultyFilter)

		buf := new(bytes.Buffer)
		err := resultsTmpl.Execute(buf, map[string]interface{}{"Courses": suggestions})
		if err != nil {
			logrus.Fatal(err)
		}

		return c.Stream(http.StatusOK, "text/html; charset=utf-8", buf)
	})

	e.Static("/vendor", "html_file/vendor")
	e.Static("/css", "html_file/css")
	e.Static("/fonts", "html_file/fonts")
	e.Static("/images", "html_file/images")
	e.Static("/js", "html_file/js")

	e.GET("/", func(c echo.Context) error {
		return c.File("./html_file/index.html")
	})

	e.Logger.Fatal(e.Start(":1234"))
}
