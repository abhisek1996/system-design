package oopsclient

import (
	lld "lld/oops_basic"
)

/*
	Encapsulation

- Go doesn't have access modifiers like private, public, or protected.
Instead, it uses the convention of starting a variable name with a lowercase letter to indicate that it is unexported (private)
and starting it with an uppercase letter to indicate that it is exported (public).
*/
func Encapsulation() {
	customer := lld.Customer{

		// name: "John Doe",   ---  cannot use unexported field
		Id: 12345,
	}

	println("Customer Name:", customer.GetName())
	println("Customer ID:", customer.GetId())
}
