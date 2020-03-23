package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const txtPath = "/txt/"

type fileHandler func(w http.ResponseWriter, r *http.Request) error

//在一个地方进行错误处理
func errWrapper(h fileHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			code := http.StatusOK
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			http.Error(w, http.StatusText(code), code)
		}
	}
}

func GetFile(w http.ResponseWriter, r *http.Request) error {

	path := r.URL.Path[len(txtPath):]
	file, err := os.Open(path)
	if err != nil {
		// panic(err)
		// fmt.Fprint(w, "%s\n", err.Error())
		return err
	}

	defer file.Close()

	content, err := ioutil.ReadAll(file)

	if err != nil {
		// panic(err)
		// fmt.Fprint(w, "%s\n", err.Error())
		return err
	}

	w.Write(content)
	return err

}

func main() {

	//route
	http.HandleFunc(txtPath, errWrapper(GetFile))

	// Start Servet
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal("ListenAndServer:   ", err)
	}
}
