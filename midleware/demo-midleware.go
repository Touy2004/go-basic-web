package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Course struct {
	CourseId   int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Instructor string  `json:"instructor"`
}

var CourseList = []Course{}

func init() {
	CourseJSON := `[
		{
			"id": 1,
			"name": "Learn Go",
			"price": 100.00,
			"instructor": "John Doe"
		},
		{
			"id": 2,
			"name": "Learn C++",
			"price": 200.00,
			"instructor": "Jane Doe"
		},
		{
			"id": 3,
			"name": "Learn C#",
			"price": 300.00,
			"instructor": "John Smith"
		} 
	]`

	err := json.Unmarshal([]byte(CourseJSON), &CourseList)

	if err != nil {
		log.Fatal(err)
	}
}

func getNextId() int {
	highestId := -1

	for _, course := range CourseList {
		if course.CourseId > highestId {
			highestId = course.CourseId
		}
	}

	return highestId + 1
}

func findId(id int) (*Course, int) {
	for _, course := range CourseList {
		if course.CourseId == id {
			return &course, id
		}
	}
	return nil, 0
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "course/")
	fmt.Println(urlPathSegment)
	id, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	course, listItemIndex := findId(id)
	if course == nil {
		http.Error(w, fmt.Sprintf("Course not found %d", id), http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:

		CourseJSON, err := json.Marshal(course)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(CourseJSON)

	case http.MethodPut:
		var updateCourse Course
		bytesBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		err = json.Unmarshal(bytesBody, &updateCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if updateCourse.CourseId != id {
			w.WriteHeader(http.StatusBadRequest)
		}
		course = &updateCourse
		CourseList[listItemIndex] = *course
		w.WriteHeader(http.StatusOK)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
func coursesHandler(w http.ResponseWriter, r *http.Request) {
	CouresJSON, err := json.Marshal(CourseList)

	switch r.Method {
	case http.MethodGet:
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(CouresJSON)
	case http.MethodPost:
		var newCourse Course
		bytesBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bytesBody, &newCourse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if newCourse.CourseId != 0 {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newCourse.CourseId = getNextId()
		CourseList = append(CourseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func midlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Before handler")
		handler.ServeHTTP(w, r)
		fmt.Println("After handler")
	})
}

func main() {
	couresItemHandler := http.HandlerFunc(courseHandler)
	CourseListHandler := http.HandlerFunc(coursesHandler)
	http.Handle("/course/", midlewareHandler(couresItemHandler))
	http.Handle("/courses", midlewareHandler(CourseListHandler))
	http.ListenAndServe(":8080", nil)
}
