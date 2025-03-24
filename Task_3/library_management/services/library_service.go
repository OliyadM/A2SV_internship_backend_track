package services

import (
	"errors"
	"sync"
	"time"
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
}

type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member
	Reservations map[int]int
	mu sync.Mutex
}

func (L *Library) AddBook(book models.Book) {
	L.Books[book.ID] = book
}

func (L *Library) RemoveBook(bookID int) error {
	if _, exists := L.Books[bookID]; !exists {
		return errors.New("book not found")
	}
	delete(L.Books, bookID)
	return nil
}


func (L *Library) BorrowBook(bookID int, memberID int) error {

	book, bookExists := L.Books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}


	member, memberExists := L.Members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	
	if book.Status != "Available" {
		return errors.New("book is already borrowed")
	}

	book.Status = "Borrowed"
	L.Books[bookID] = book 


	member.BorrowedBooks = append(member.BorrowedBooks, book)
	L.Members[memberID] = member 

	return nil
}


func (L *Library) ReturnBook(bookID int, memberID int) error {
	book, bookExists := L.Books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}
	member, memberExists := L.Members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}
	var bookIndex int
	found := false
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			bookIndex = i
			found = true
			break
		}
	}
	if !found {
		return errors.New("member has not borrowed this book")
	}
	member.BorrowedBooks = append(member.BorrowedBooks[:bookIndex], member.BorrowedBooks[bookIndex+1:]...)
	L.Members[memberID] = member
	book.Status = "Available"
	L.Books[bookID] = book
	return nil
}

func (L *Library) ListAvailableBooks() []models.Book {
	var availableBooks []models.Book
	for _, book := range L.Books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (L *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	member, exists := L.Members[memberID]
	if !exists {
		return nil, errors.New("member not found")
	}
	return member.BorrowedBooks, nil
}

func (L *Library) ReserveBook(bookID int, memberID int) error {
	L.mu.Lock()
	defer L.mu.Unlock()

	book, exists := L.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}


	if _, reserved := L.Reservations[bookID]; reserved {
		return errors.New("book is already reserved")
	}

	L.Reservations[bookID] = memberID

	
	go func() {
		time.Sleep(5 * time.Second)

		L.mu.Lock()
		defer L.mu.Unlock()

		
		if L.Reservations[bookID] == memberID && book.Status == "Available" {
			delete(L.Reservations, bookID)
		}
	}()

	return nil
}
