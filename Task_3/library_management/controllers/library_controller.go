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
		fmt.Println("7. reserve books for a member")
		fmt.Println("8. Exit")
		fmt.Print("Choose an operation: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			addbook()
		case 2:
			removebook()
		case 3:
			borrowbook()
		case 4:
			returnbook()
		case 5:
			listavailablebooks()
		case 6:
			listborrowedbooks()
		case 7:
			reservebook()
		case 8:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please choose from the listed operations.")
		}
	}
}

func addbook() {
	var id int
	var title, author, status string

	fmt.Print("Enter book ID: ")
	fmt.Scan(&id)
	fmt.Print("Enter book title: ")
	fmt.Scan(&title)
	fmt.Print("Enter book author: ")
	fmt.Scan(&author)
	fmt.Print("Enter book status (Available/Borrowed): ")
	fmt.Scan(&status)

	newBook := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: status,
	}

	library.AddBook(newBook)
	fmt.Println("Book added successfully!")
}

func removebook() {
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

func borrowbook() {
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

func returnbook() {
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

func listavailablebooks() {
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

func listborrowedbooks() {
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

func reservebook() {
	var memberID, bookID int

	fmt.Print("Enter Member ID: ")
	fmt.Scan(&memberID)
	fmt.Print("Enter Book ID: ")
	fmt.Scan(&bookID)

	err := library.ReserveBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book reserved successfully!")
	}
}
