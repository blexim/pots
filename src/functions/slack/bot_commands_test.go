package main

import "fmt"

func ExampleParseBuyin() {
	fmt.Println(ParseBuyin("alice", "@pokerbot in £13.70"))
	// Output: &{alice 1370}
}

func ExampleParseBuyin2() {
	fmt.Println(ParseBuyin("alice", "@pokerbot in £13"))
	// Output: &{alice 1300}
}

func ExampleParseBuyin3() {
	fmt.Println(ParseBuyin("alice", "@pokerbot in 20"))
	// Output: &{alice 2000}
}
