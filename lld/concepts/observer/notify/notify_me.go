package notify

func StockManager() {
	// Create Electronics Item
	iphone := NewElectronicsItem("iPhone 12", 100)
	iphone.AddObserver(NewEmailAlert("abhishek.pradhan@abc.com"))
	iphone.AddItem(1)

	// Create Grocery Item
	milk := NewGroceryItem("Milk")
	milk.AddObserver(NewMobileAlert("1234567890"))
	milk.AddItem(2)
}
