package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Room struct {
	Name           string
	Code           string
	TextIntro      string
	Text           func() string
	things         map[string]string
	availablePath  map[string]string
	active_objects map[string]string
	doorOpen       bool
}

type Person struct {
	inventory []string
	Room
}

var allRooms map[string]Room
var equipment map[string]string
var commands map[string]string
var Vasily Person

func main() {
	initGame()
	in := bufio.NewScanner(os.Stdin)
Loop:
	for in.Scan() {
		command := in.Text()
		fmt.Fprintln(os.Stdout, handleCommand(command))
		fmt.Fprintln(os.Stdout, allRooms["комната"].things)
		if Vasily.Room.Name == "улица" {
			break Loop
		}
	}
}

func initGame() {
	var kitchen Room = Room{
		Name:      "коридор",
		Code:      "kitchen",
		TextIntro: "кухня, ничего интересного. можно пройти - коридор",
		Text: func() string {
			_, synopsis := allRooms["комната"].things["конспекты"]
			_, bag := allRooms["комната"].things["рюкзак"]
			_, keys := allRooms["комната"].things["ключи"]
			switch {
			case synopsis || bag || keys:
				return "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"
			default:
				return "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"
			}
		},
		things: map[string]string{
			"чай": "tea",
		},
		availablePath: map[string]string{
			"коридор": "hollway",
		},
	}

	var hollway Room = Room{
		Name:      "коридор",
		Code:      "hollway",
		TextIntro: "ничего интересного. можно пройти - кухня, комната, улица",
		Text:      func() string { return "ничего интересного" },
		availablePath: map[string]string{
			"кухня":   "kitchen",
			"комната": "Room",
			"улица":   "street",
		},
		active_objects: map[string]string{
			"дверь": "hollway",
		},
	}

	var room Room = Room{
		Name:      "комната",
		Code:      "room",
		TextIntro: "ты в своей комнате. можно пройти - коридор",
		Text: func() string {
			_, synopsis := allRooms["комната"].things["конспекты"]
			_, bag := allRooms["комната"].things["рюкзак"]
			_, keys := allRooms["комната"].things["ключи"]
			switch {
			case synopsis && bag && keys:
				return "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор"
			case synopsis && bag == false && keys:
				return "на столе: ключи, конспекты. можно пройти - коридор"
			case synopsis && bag == false && keys == false:
				return "на столе: конспекты. можно пройти - коридор"
			case synopsis == false && bag == false && keys:
				return "на столе: ключи. можно пройти - коридор"
			default:
				return "пустая комната. можно пройти - коридор"
			}
		},
		things: map[string]string{
			"конспекты": "synopsis",
			"ключи":     "keys",
			"рюкзак":    "bag",
		},
		availablePath: map[string]string{
			"коридор": "hollway",
			"кухня":   "kitchen",
			"улица":   "street",
		},
	}

	var street Room = Room{
		Name:      "улица",
		Code:      "street",
		TextIntro: "на улице весна. можно пройти - домой",
		Text:      func() string { return "на улице весна" },
	}

	allRooms = map[string]Room{
		"коридор": hollway,
		"кухня":   kitchen,
		"комната": room,
		"улица":   street,
	}
	Vasily = Person{
		Room: allRooms["кухня"],
	}
	equipment = map[string]string{
		"рюкзак": "bag",
	}
	commands = map[string]string{
		"осмотреться": "lookAround",
		"идти":        "move",
		"надеть":      "wear",
		"взять":       "take",
		"применить":   "use",
	}
}

func handleCommand(Text string) string {
	command := strings.Split(Text, " ")
	if len(command) == 0 {
		return "Введите команду!"
	}
	switch commands[command[0]] {
	case "lookAround":
		return lookAround()
	case "move":
		return move(command)
	case "wear":
		return wear(command)
	case "take":
		return take(command)
	case "use":
		return use(command)
	default:
		return "неизвестная команда"
	}
}

func lookAround() string {
	return Vasily.Room.Text()
}

func move(command []string) string {
	if len(command) < 2 {
		return "Введите аргументы!!!"
	}
	if _, roomExist := allRooms[command[1]]; !roomExist {
		return "нет пути в " + command[1]
	}
	if _, availablePathExist := Vasily.availablePath[command[1]]; !availablePathExist {
		return "нет пути в " + command[1]
	}
	if Vasily.availablePath[command[1]] == "street" {
		if Vasily.Room.doorOpen {
			Vasily.Room = allRooms[command[1]]
		} else {
			return "дверь закрыта"
		}
	} else {
		Vasily.Room = allRooms[command[1]]
	}

	return Vasily.TextIntro
}

func wear(command []string) string {
	if len(command) < 2 {
		return "Введите аргументы!!!"
	}
	if _, equipmentExist := equipment[command[1]]; equipmentExist {
		Vasily.inventory = append(Vasily.inventory, command[1])
		delete(allRooms[Vasily.Room.Name].things, command[1])
		delete(Vasily.things, command[1])
	} else {
		return "нет такого"
	}
	return "вы надели: " + command[1]
}

func take(command []string) string {
	if len(command) < 2 {
		return "Введите аргументы!!!"
	}
	if _, knownThing := Vasily.things[command[1]]; !knownThing {
		return "нет такого"
	}
	if len(Vasily.inventory) > 0 {
		Vasily.inventory = append(Vasily.inventory, command[1])
		delete(allRooms[Vasily.Room.Name].things, command[1])
		delete(Vasily.things, command[1])
	} else {
		return "некуда класть"
	}
	return "предмет добавлен в инвентарь: " + command[1]
}

func use(command []string) string {
	if len(command) < 3 {
		return "Введите аргументы!!!"
	}
	firstObject := false
Loop:
	for i := 1; i < len(Vasily.inventory); i++ {
		if Vasily.inventory[i] == command[1] {
			firstObject = true
			break Loop
		}
	}
	if !firstObject {
		return "нет предмета в инвентаре - " + command[1]
	}
	if _, secondObject := Vasily.active_objects[command[2]]; secondObject {
		Vasily.Room.doorOpen = true
		return "дверь открыта"
	} else {
		return "не к чему применить"
	}
}
