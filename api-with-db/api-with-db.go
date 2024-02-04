package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var (
	userDB string = "root"
	passDB string = ""
	hostDB string = "localhost"
	portDB string = "3306"
	dbName string = "coursesdb"
)

const basePath = "/api"
const coursePath = "courses"

type Course struct {
	CourseId int     `json:"courseid"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"imageurl"`
}

func SetupDB() {
	var err error

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userDB, passDB, hostDB, portDB, dbName))

	if err != nil {
		log.Println("failed to connect database")
	} else {
		log.Println("connected to database")
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func getCourseList() ([]Course, error) {
	cxt, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result, err := db.QueryContext(cxt, "SELECT courseid, coursename, price, image_url FROM courseonline")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer result.Close()

	courses := make([]Course, 0)
	for result.Next() {
		var course Course
		err := result.Scan(&course.CourseId, &course.Name, &course.Price, &course.ImageURL)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

func insertProduct(course Course) (int, error) {
	cxt, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	result, err := db.ExecContext(cxt, "INSERT INTO courseonline(courseid, coursename, price, image_url) VALUES (?, ?, ?, ?)", course.CourseId, course.Name, course.Price, course.ImageURL)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertId), nil
}

func getCourse(courseId int) (*Course, error) {
	cxt, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	row := db.QueryRowContext(cxt, "SELECT courseid, coursename, price, image_url FROM courseonline WHERE courseid = ?", courseId)

	course := &Course{}
	err := row.Scan(&course.CourseId, &course.Name, &course.Price, &course.ImageURL)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil
}

func deleteCourse(courseId int) error {
	cxt, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := db.ExecContext(cxt, "DELETE FROM courseonline WHERE courseid = ?", courseId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func updateCourse(course Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := db.ExecContext(ctx, "UPDATE courseonline SET courseid = ?, coursename = ?, price = ?, image_url = ? WHERE courseid = ?", course.CourseId, course.Name, course.Price, course.ImageURL, course.CourseId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func handleCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		courseList, err := getCourseList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		courseJSON, err := json.Marshal(courseList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(courseJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		var course Course
		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		CourseId, err := insertProduct(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"courseid": %d}`, CourseId)))

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleCourse(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))
	if len(urlPathSegment[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	courseId, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		course, err := getCourse(courseId)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		courseJSON, err := json.Marshal(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(courseJSON)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		err := deleteCourse(courseId)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		var course Course
		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = updateCourse(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	}
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		handler.ServeHTTP(w, r)
	})
}

func SetRoutes(apiBasePath string) {
	courseHandler := http.HandlerFunc(handleCourse)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler))
	coursesHandler := http.HandlerFunc(handleCourses)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler))
}

func main() {
	SetupDB()
	SetRoutes(basePath)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
