package taskss

import "fmt"

func Word_count( s string)  map[string]int {
	var freq_dict = make(map[string] int)

	for _, letter := range s {
		freq_dict[string(letter)] += 1
	}
	return freq_dict

}

func Palindrom(s string) bool{
	var n = len(s)
	var left = 0
	var right = n-1

	for left<right{
		if s[left] != s[right]{
			return false
		}
		left+=1
		right-=1
	}
	return true
}


func Task_2(){
var word string

fmt.Println("Enter a word")
fmt.Scan(&word)

fmt.Println(Word_count(word))

if Palindrom(word) {
	fmt.Printf("%v is a palindromic word \n" , word)
} else {
	fmt.Printf("%v is a not a palindromic word \n" , word)
}


}


