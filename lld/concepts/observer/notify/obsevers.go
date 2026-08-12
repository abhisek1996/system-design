package notify

import "fmt"

// --- Interfaces ---

type Observer interface {
	Update(e Event)
}

// --- Concrete Observers ---

type EmailAlert struct {
	EmailId string
}

func NewEmailAlert(emailId string) *EmailAlert {
	return &EmailAlert{EmailId: emailId}
}

func (e *EmailAlert) Update(event Event) {
	// Extract data safely or fall back to defaults
	var name string
	// Check "model" first, then "name"
	if val, ok := event.Data["model"]; ok {
		name = val.(string)
	} else if val, ok := event.Data["name"]; ok {
		name = val.(string)
	} else {
		name = "Unknown Item"
	}

	var stock int
	if val, ok := event.Data["stock"]; ok {
		stock = val.(int)
	}

	fmt.Printf("Email to %s: [%s] %s (Stock: %d)\n", e.EmailId, event.Type, name, stock)
}

type MobileAlert struct {
	MobileNumber string
}

func NewMobileAlert(mobileNumber string) *MobileAlert {
	return &MobileAlert{MobileNumber: mobileNumber}
}

func (m *MobileAlert) Update(event Event) {
	// SMS just sends the message
	fmt.Printf("SMS to %s: %s\n", m.MobileNumber, event.Message)
}
