package main

import (
	"fmt"
	"kiber/config"
	"kiber/transport/customers"
	"kiber/transport/groups"
	"kiber/transport/kiberons"
	"kiber/transport/teacher"
	"kiber/transport/token"
)

func main() {
	config, err := config.LoadConfig()

	if err != nil {
		fmt.Println("Ошибка конфига !")
		return
	}
	err = token.GetToken(config)

	if err != nil {
		fmt.Println("%v\n", err)
		return
	}

	fmt.Println("Привет!")

	fmt.Println("\nСделаем выборку по тьюторам:")

	for _, v := range config.Names_tutors {
		fmt.Println(v)
	}

	err = teacher.GetAllIdsTeachers(config)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = groups.GetAllGroups(config)

	fmt.Println("Айди детей, которые заслужили кибероны:")

	data, err := customers.GetCustomers(config)

	for _, val := range data {
		fmt.Println(val)
	}

	fmt.Println("Начинаю выставлять")

	kiberons.AddBonusesToUsers(config, data)
}
