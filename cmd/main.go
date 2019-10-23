package main

import (
	"github.com/siddhantcse14/go-promise/promise"
	"log"
	"time"
)

func main() {

	// Running Example
	response := promise.New(func(resolve, reject func(interface{})) {
		time.Sleep(2 * time.Second)
		resolve(2)
	}).Then(func(value interface{}) interface{} {
		log.Println("In first Then")
		return value.(int) + 4
	}).Then(func(value interface{}) interface{} {
		log.Println("In second Then")
		panic("Panic Occurred due to Testing")
		return value.(int) + 20
	}).Catch(func(reason interface{}) interface{} {
		log.Println("Panic Catched, reason: ",reason.(string))
		return reason
	}).Finally(func(result interface{}) interface{}{
		log.Println("In Finally")
		return result
	}).Await()

	log.Println(response)
}