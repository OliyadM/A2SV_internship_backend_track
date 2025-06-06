package controllers

import (
	"fmt"
	"library_management/models"
	"library_management/services"
)

var library = services.Library{
	Books:   make(map[int]models.Book),
	Members: make(map[int]models.Member),
}

func UserInput() {
	for {
		fmt.Println("\nLibrary Management System")
		fmt.Println("1. Add a new book")
		fmt.Println("2. Remove a book")
		fmt.Println("3. Borrow a book")
		fmt.Println("4. Return a book")
		fmt.Println("5. List available books")
		fmt.Println("6. List borrowed books by a member")
		fmt.Println("7. Exit")
		fmt.Print("Choose an operation: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			addBook()
		case 2:
			removeBook()
		case 3:
			borrowBook()
		case 4:
			returnBook()
		case 5:
			listAvailableBooks()
		case 6:
			listBorrowedBooks()
		case 7:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please choose from the listed operations.")
		}
	}
}

func addBook() {
	var id int
	var title, author string

	fmt.Print("Enter book ID: ")
	fmt.Scan(&id)
	fmt.Print("Enter book title: ")
	fmt.Scan(&title)
	fmt.Print("Enter book author: ")
	fmt.Scan(&author)

	newBook := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: "Available",
	}

	library.AddBook(newBook)
	fmt.Println("Book added successfully!")
}

func removeBook() {
	var id int
	fmt.Print("Enter book ID to remove: ")
	fmt.Scan(&id)

	err := library.RemoveBook(id)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book removed successfully!")
	}
}

func borrowBook() {
	var bookID, memberID int
	fmt.Print("Enter book ID: ")
	fmt.Scan(&bookID)
	fmt.Print("Enter member ID: ")
	fmt.Scan(&memberID)

	err := library.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book borrowed successfully!")
	}
}

func returnBook() {
	var bookID, memberID int
	fmt.Print("Enter book ID: ")
	fmt.Scan(&bookID)
	fmt.Print("Enter member ID: ")
	fmt.Scan(&memberID)

	err := library.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book returned successfully!")
	}
}

func listAvailableBooks() {
	books := library.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}
	fmt.Println("Available Books:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func listBorrowedBooks() {
	var memberID int
	fmt.Print("Enter member ID: ")
	fmt.Scan(&memberID)

	books, err := library.ListBorrowedBooks(memberID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Borrowed Books:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s\n", book.ID, book.Title)
	}
}
