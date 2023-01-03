package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var availableOperations = "+-/*"

const maxArgumentsCount = 2
const maxInput = 10
const minInput = 1
const maxInputRoman = "X"
const minInputRoman = "I"

var romanNumbers = map[int]string{
	1:   "I",
	2:   "II",
	3:   "III",
	4:   "IV",
	5:   "V",
	6:   "VI",
	7:   "VII",
	8:   "VIII",
	9:   "IX",
	10:  "X",
	40:  "XL",
	50:  "L",
	90:  "XC",
	100: "C",
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Введите математическое выражение с числами от 1 до 10 в римском или арабском формате")

		input, err := reader.ReadString('\n')

		if input == "exit" {
			break
		}

		output, err := calculate(input)

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(output)
	}
}

func calculate(input string) (output string, err error) {
	input = strings.TrimSpace(input)
	input = strings.Replace(input, " ", "", -1)

	arguments, mathOperator, err := getArguments(input)

	if err != nil {
		return input, err
	}

	isRoman, err := checkArgumentsFormat(arguments)

	if err != nil {
		return input, err
	}

	argumentsArabic, err := convertInt(arguments, isRoman)

	if err != nil {
		return input, err
	}

	finalValue := doMath(argumentsArabic, mathOperator)

	if isRoman && (finalValue <= 0) {
		err = fmt.Errorf("результат выражения равен %d и невозможен для вывода в римских числах", finalValue)

		return input, err
	}

	if isRoman {
		output = convertArabicRome(finalValue)
		return output, err
	}

	output = strconv.Itoa(finalValue)

	return output, err
}

func getArguments(input string) (arguments []string, mathOperator string, err error) {
	for index := range availableOperations {
		if strings.Contains(input, string(availableOperations[index])) {
			mathOperator = string(availableOperations[index])
		}
	}

	if mathOperator == "" {
		err = errors.New("выражение не содержит математического оператора")
		return arguments, mathOperator, err
	}

	arguments = strings.Split(input, mathOperator)

	if len(arguments) > maxArgumentsCount {
		err = fmt.Errorf("количество аргументов больше %d", maxArgumentsCount)
		return arguments, mathOperator, err
	}

	return arguments, mathOperator, nil
}

func checkArgumentsFormat(arguments []string) (isRoman bool, err error) {
	isRomanFirst := checkRoman(arguments[0])
	isRomanSecond := checkRoman(arguments[1])

	if isRomanFirst != isRomanSecond {
		err = fmt.Errorf("введены числа из разных систем или же одно из чисел не входит в диапазон от 1 до 10")
		return isRoman, err
	}

	if isRomanFirst && isRomanSecond {
		return true, nil
	}

	return false, nil
}

func checkRoman(argument string) (isRoman bool) {
	for _, value := range romanNumbers {
		if strings.ToLower(argument) == strings.ToLower(value) {
			return true
		}
	}

	return isRoman
}

func convertInt(arguments []string, isRoman bool) (argumentsInt []int, err error) {
	if isRoman {
		return convertArrayRomeArabic(arguments)
	}

	for i := 0; i < maxArgumentsCount; i++ {
		temporalArgument, err := strconv.Atoi(arguments[i])

		if err != nil {
			err = errors.New("один или оба из аргументов не являются числомами, с которыми возможна операция")
			return argumentsInt, err
		}

		if temporalArgument > 10 || temporalArgument < 1 {
			err = fmt.Errorf("операция возможна только с числами от %d до %d", minInput, maxInput)
			return argumentsInt, err
		}
		argumentsInt = append(argumentsInt, temporalArgument)
	}

	return argumentsInt, err
}

func convertArrayRomeArabic(arguments []string) (argumentsInt []int, err error) {
	var temporaryArgument int

	for i := 0; i < maxArgumentsCount; i++ {
		temporaryArgument, err = convertRomeArabic(arguments[i])
		argumentsInt = append(argumentsInt, temporaryArgument)
	}
	return argumentsInt, err
}

func convertRomeArabic(argument string) (argumentInt int, err error) {
	for key, value := range romanNumbers {
		if strings.ToLower(value) == strings.ToLower(argument) && (key >= minInput && key <= maxInput) {
			return key, err
		}
	}

	return 0, err
}

func doMath(argumentsArabic []int, mathOperator string) int {
	switch mathOperator {
	case "+":
		return argumentsArabic[0] + argumentsArabic[1]

	case "-":
		return argumentsArabic[0] - argumentsArabic[1]

	case "*":
		return argumentsArabic[0] * argumentsArabic[1]

	case "/":
		return argumentsArabic[0] / argumentsArabic[1]
	}

	return 0
}

func convertArabicRome(arabicNumber int) (romanNumber string) {
	var arguments []int
	divisor := 10

	for arabicNumber > 0 {
		arguments = append(arguments, arabicNumber%divisor)
		arabicNumber -= arabicNumber % divisor
		divisor *= 10
	}

	for i := len(arguments) - 1; i >= 0; i-- {
		romanNumber += getRomeNumber(arguments[i])
	}

	return
}

func getRomeNumber(argument int) (romeNumber string) {
	var temporaryValue string
	maxKey := 0

	for argument > 0 {
		for key, value := range romanNumbers {
			if maxKey < key && key <= argument {
				maxKey = key
				temporaryValue = value
			}
		}

		argument -= maxKey
		romeNumber += temporaryValue
		maxKey = 0
	}

	return romeNumber
}
