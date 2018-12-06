package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Operation struct {
	Description 	string 	`json:"description"`
	Result			float64	`json:"result"`
}

type Person struct {
	Name 			string 	`json:"name"`
	Email			string	`json:"email"`
}

const msgNotFound= "{ \"Message\": \"Not Found.\" }"
const msgNotImplemented = "{ \"Message\": \"Not implemented.\" }"
const msgError = "{ \"Message\": \"Sometimes things don't work as expected...\" }"

func main() {
	var port = 8088

	/* Router */
	router := mux.NewRouter()
	// Catch 404
	router.NotFoundHandler = http.HandlerFunc(NotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(NotAllowed)

	/* Paths */

	// Base Path [localhost:port/]
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type","application/json")
		response.WriteHeader(http.StatusOK)
		response.Write(json.RawMessage("{ \"Message\": \"REST API in GO!!!\" }"))
	})

	// Create SubRouter [localhost:port/calc]
	subRouterCalc := router.PathPrefix("/calc").Subrouter()
	// Map paths on subRouter to methods
	subRouterCalc.HandleFunc("/sum/{num1}/{num2}", CalcSum).Methods("GET") // [localhost:port/calc/sum/{num1}/{num2}]
	subRouterCalc.HandleFunc("/sub/{num1}/{num2}", DoItYourself).Methods("GET")
	subRouterCalc.HandleFunc("/mul/{num1}/{num2}", DoItYourself).Methods("GET")
	subRouterCalc.HandleFunc("/div/{num1}/{num2}", DoItYourself).Methods("GET")

	// Another SubRouter [localhost:port/calcall]
	subRouterCalcAll := router.PathPrefix("/calcall").Subrouter()
	subRouterCalcAll.HandleFunc("/{operator}/{num1}/{num2}", Calc).Methods("GET")	 // [localhost:port/calcall/sum/{num1}/{num2}]

	// POST
	router.HandleFunc("/validemail", PostTest).Methods("POST")

	/* Start Server */
	println("Running http server @ localhost:" + strconv.Itoa(port))
	if err := http.ListenAndServe(":" + strconv.Itoa(port), router); err != nil {
		// Log & Exit
		log.Fatal(err)
	}
}

// 404 Handler
// Handles invalid paths
func NotFound(response http.ResponseWriter, request *http.Request) {
	// Print url to console
	log.Println("Not Found: " + request.RequestURI)
	// Response
	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusNotFound)
	response.Write(json.RawMessage(msgNotFound))
}

// 405
// Handles not allowed methods (GET/POST/PUT/DELETE)
func NotAllowed(response http.ResponseWriter, request *http.Request) {
	// Print method and url to console
	log.Println("Not allowed: " + request.Method + "@" + request.RequestURI)
	// Response
	response.Header().Set("Content-Type","application/json")
	response.WriteHeader(http.StatusMethodNotAllowed)
	response.Write(json.RawMessage(http.ErrBodyNotAllowed.Error()))
}

// Test GET http://localhost:8080/calc/sum/2/4
func CalcSum(resp http.ResponseWriter, req *http.Request) {
	// Params
	vars := mux.Vars(req)
	operation := Operation{}
	operation.Description = fmt.Sprintf("%s + %s", vars["num1"], vars["num2"])
	// Get int values
	num1, err := strconv.ParseFloat(vars["num1"], 64)
	num2, err := strconv.ParseFloat(vars["num2"], 64)
	if err == nil {
		operation.Result = num1 + num2
		output, err := json.Marshal(operation)
		// Response
		resp.Header().Set("Content-Type","application/json")
		if err == nil {
			// Response status
			resp.WriteHeader(http.StatusOK)
			resp.Write(output)
		} else {
			panic("Something didn't work...")
		}
	} else {
		log.Println(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write(json.RawMessage(msgError))
	}
}

//region TODO
// Sub / Mul / Div
func DoItYourself(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type","application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(json.RawMessage(msgNotImplemented))
}
//endregion

// Calc All
func Calc(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type","application/json")
	// Params
	vars := mux.Vars(req)
	re := regexp.MustCompile("(sum|sub|div|mul)")
	var operator = vars["operator"]
	num1, err1 := strconv.ParseFloat(vars["num1"],64)
	num2, err2 := strconv.ParseFloat(vars["num2"],64)
	if re.MatchString(operator) && err1 == nil && err2 == nil {
		operation := Operation{}
		var temp string
		switch operator {
			case "sum":
				operation.Result = num1 + num2
				temp = "+"
			case "sub":
				operation.Result = num1 - num2
				temp = "-"
			case "mul":
				operation.Result = num1 * num2
				temp = "*"
			case "div":
				if num2 != 0 {
					operation.Result = num1 / num2
					temp = "/"
				} else {
					resp.WriteHeader(http.StatusBadRequest)
					resp.Write(json.RawMessage("{ \"Message\": \"Can't divide by 0.\" }"))
					return
				}
		}
		operation.Description = fmt.Sprintf("%s %s %s", vars["num1"], temp, vars["num2"])
		var output, _ = json.Marshal(operation)
		resp.WriteHeader(http.StatusOK)
		resp.Write(output)
	} else {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write(json.RawMessage(msgError))
	}
}

/* Test POST
	{
		"name": "John Doe",
		"email": "johndoe@doe.com"
	}
*/
func PostTest(resp http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	person := Person{}
	err := decoder.Decode(&person)
	resp.Header().Set("Content-Type","application/json")
	if err != nil {
		log.Println(err)
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(json.RawMessage(msgError))
	}
	resp.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("The email %s of %s is %s", person.Email, person.Name, IsEmailValid(person.Email))
	resp.Write(json.RawMessage("{ \"Message\": \""+ message + "\" }"))
}

// Check if email is valid
func IsEmailValid(email string) string {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(email) {
	 	return "Valid"
	 } else {
	 	return "Invalid"
	}
}