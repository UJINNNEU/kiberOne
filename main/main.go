package main

import (
	"fmt"
	"kiber/config"
	"kiber/transport/customers"
	"kiber/transport/groups"
	"kiber/transport/teacher"
)

func main() {
	config, _ := config.LoadConfig()
	fmt.Println("Привет!")

	fmt.Println("\nСделаем выборку по тьюторам:")

	for _, v := range config.Names_tutors {
		fmt.Println(v)
	}

	//Получение ID у тьюторов
	err := teacher.GetAllIdsTeachers(config)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = groups.GetAllGroups(config)

	fmt.Println("Good kid:")

	data, err := customers.GetCustomers(config)

	for _, val := range data {
		fmt.Println(val)
	}

}
