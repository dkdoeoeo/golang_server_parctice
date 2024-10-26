package main

import "post-platform/router"

func main() {
	e := router.Router()
	e.Run(":8080")
}
