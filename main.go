package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// maxValueFrom возвращает верхнюю границу диапазона генерации
// в зависимости от уровня сложности: easy=10, hard=1000, default=100.
func maxValueFrom(difficulty string) uint32 {
	switch difficulty {
	case "easy":
		return 10
	case "hard":
		return 1000
	default:
		return 100
	}
}

// guessLoop читает угадывания из stdin и сравнивает их с secretNumber.
// Возвращает true, если число угадано.
func guessLoop(secretNumber uint32) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Please input your guess.")

		line, err := reader.ReadString('\n')
		if err != nil {
			// EOF или другая ошибка чтения — выходим из цикла.
			return false
		}

		guess, err := strconv.ParseUint(strings.TrimSpace(line), 10, 32)
		if err != nil {
			continue
		}
		guessValue := uint32(guess)

		fmt.Printf("You guessed: %d\n", guessValue)

		switch {
		case guessValue < secretNumber:
			fmt.Println("Too small!")
		case guessValue > secretNumber:
			fmt.Println("Too big!")
		default:
			fmt.Println("You win!")
			return true
		}
	}
}

func main() {
	fmt.Println("Guess the number!")

	difficulty := "normal"
	if len(os.Args) > 1 {
		difficulty = os.Args[1]
	}
	maxValue := maxValueFrom(difficulty)
	secretNumber := uint32(rand.Intn(int(maxValue))) + 1

	fmt.Printf("The secret number is: %d (difficulty: %s, max: %d)\n", secretNumber, difficulty, maxValue)

	guessLoop(secretNumber)
}
