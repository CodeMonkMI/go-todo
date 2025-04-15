package utility

import (
	"fmt"
	"log"
	"net/http"

	"github.com/thedevsaddam/renderer"
)

type Payload struct {
	rnd *renderer.Render
}

func ResponseError(w http.ResponseWriter, err error, msg string) {
	log.Println(err)
	rnd := renderer.New()
	rnd.JSON(w, http.StatusBadRequest, renderer.M{
		"message": msg,
		"err":     err,
	})
	return
}

func CheckErr(e error, customMsg string) {
	if e != nil {
		fmt.Println(customMsg)
		log.Fatal(e)
	}
}
