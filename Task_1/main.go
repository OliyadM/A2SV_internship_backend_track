package main

import "fmt"

func main() {
	task1()
}

func displayScore(scoreDict map[string]int) {
	for subject, score := range scoreDict {
		fmt.Printf("Your score for %v is %v \n", subject, score)
	}
}

func average(scoreDict map[string]int, n int) {
	var total int
	for _, score := range scoreDict {
		total += score
	}
	ave := float64(total) / float64(n)
	fmt.Printf("Your average score is %.2f \n", ave)
}

func task1() {
	var name string
	var n int
	var subject string
	var score int
	scoreDict := make(map[string]int)

	fmt.Println("Enter your name: ")
	fmt.Scanln(&name)

	fmt.Println("Enter the number of subjects you take:")
	fmt.Scanln(&n)

	for i := 0; i < n; i++ {
		fmt.Printf("Enter subject %v: ", i+1)
		fmt.Scanln(&subject)

		fmt.Printf("Enter the score for %v: ", subject)
		fmt.Scanln(&score)

		scoreDict[subject] = score
	}

	displayScore(scoreDict)
	average(scoreDict, n)
}
