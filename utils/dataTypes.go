package dataTypes

type Tutorial struct {
	ID int
	Title string
	Author Author
	Comments []Comment
}

type Author struct {
	Name string
	Tutorials []int
}

type Comment struct {
	Body string
}