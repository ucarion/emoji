package emoji_test

import (
	"fmt"

	"github.com/ucarion/emoji"
)

func ExampleLookup() {
	fmt.Println(emoji.Lookup("a"))
	fmt.Println(emoji.Lookup("😎"))

	// Output:
	// { 0  } false
	// {smiling face with sunglasses 1 1.0 😎} true
}

func ExampleLookup_Status() {
	s := "☺"
	e1, _ := emoji.Lookup(s)
	e2, _ := emoji.Lookup(e1.FullyQualifiesAs)

	fmt.Println(s, e1.FullyQualifiesAs)
	fmt.Printf("%U %U\n", []rune(s), []rune(e1.FullyQualifiesAs))
	fmt.Println(e1.Status == emoji.Unqualified, e2.Status == emoji.FullyQualified)

	// Output:
	// ☺ ☺️
	// [U+263A] [U+263A U+FE0F]
	// true true
}
