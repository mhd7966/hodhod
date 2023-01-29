package repositories

import (
	"errors"

	"github.com/abr-ooo/hodhod/connections"
	"github.com/abr-ooo/hodhod/inputs"
	"github.com/abr-ooo/hodhod/models"
	// log "github.com/sirupsen/logrus"
)

// func Exist(table string, field string, value string) bool {

// 	var found bool
// 	query := "SELECT EXISTS(SELECT 1 FROM " + table + "WHERE " + field + " = $1) AS found"

// 	connections.DB.Raw(query, value).Scan(&found)

// 	return found
// }

func ExistContact(userID int, input inputs.ContactBody) bool {
	var contact models.Contact

	if result := connections.DB.Where(models.Contact{
		UserID: userID,
		Name:   input.Name,
	}).First(&contact); result.Error != nil {
		return false
	}

	return true
}

// func ExistRecord(table string, field1 string, value1 string, field2 string, value2 string) bool {

// 	var found bool
// 	query := "SELECT EXISTS(SELECT 1 FROM " + table + " WHERE " + field1 + " = $1 and " + field2 + " = $2) AS found"
// 	fmt.Println(query)

// 	connections.DB.Raw(query, value1, value2).Scan(&found)
// 	fmt.Println(found)

// 	return true
// }

func GetContact(contactID int) (*models.Contact, error) {
	var contact models.Contact

	if result := connections.DB.First(&contact, contactID); result.Error != nil {
		return nil, result.Error
	}

	return &contact, nil
}

func GetContacts(userID int) (*[]models.Contact, error) {
	var contacts []models.Contact

	if result := connections.DB.Find(&contacts, models.Contact{UserID: userID}); result.RowsAffected < 1 {
		return nil, errors.New("user doesn't have any contact")
	}

	return &contacts, nil
}

func CreateContact(userID int, input inputs.ContactBody) (*models.Contact, error) {
	contact := models.Contact{
		UserID:      userID,
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
	}
	if result := connections.DB.Create(&contact); result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

func UpdateContact(contact models.Contact, contactID int, input inputs.ContactBody, name bool, email bool, phoneNumber bool) error {
	if name {
		contact.Name = input.Name
	}
	if email {
		contact.Email = input.Email
	}
	if phoneNumber {
		contact.PhoneNumber = input.PhoneNumber
	}

	if result := connections.DB.Save(&contact); result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteContact(contact models.Contact) error {

	if result := connections.DB.Delete(&contact); result.Error != nil {
		return result.Error
	}

	return nil
}

func GetContactList(ids []int) ([]models.Contact, error) {

	var contacts []models.Contact

	if result := connections.DB.Find(&contacts, ids); result.RowsAffected != int64(len(ids)) {
		return nil, errors.New("all contacts doesn't found")
	}

	return contacts, nil
}
