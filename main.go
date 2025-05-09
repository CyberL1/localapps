package main

import (
	"localapps/cmd"
	"localapps/utils"
)

func main() {
	utils.UpdateCliConfigCache()
	cmd.Execute()
}
