package seeder

import (
	"log"

	"aro-shop/db"
	"aro-shop/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateSuperAdminIfNotExists() {
	var existing models.User

	err := db.DB.Where("email = ?", "superadmin@gmail.com").First(&existing).Error
	if err == nil {
		log.Println("[Seeder] Super Admin already exists")
		return
	} else if err != gorm.ErrRecordNotFound {
		log.Println("[Seeder] Failed to check super admin:", err)
		return
	}

	password := "supersecret"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	superAdmin := models.User{
		Name:     "Super Admin",
		Email:    "superadmin@gmail.com",
		Password: string(hashedPassword),
		Role:     models.RoleSuperAdmin,
	}

	if err := db.DB.Create(&superAdmin).Error; err != nil {
		log.Println("[Seeder] Failed to create Super Admin:", err)
	} else {
		log.Println("[Seeder] Super Admin created: superadmin@gmail.com | password:", password)
	}
}

func SeedPaymentMethods() error {
	paymentMethods := []models.PaymentMethod{
		{Name: "Cash"},
		{Name: "Credit Card"},
		{Name: "Debit Card"},
		{Name: "Bank Transfer"},
		{Name: "E-Wallet"},
		{Name: "QRIS"},
		{Name: "GoPay"},
		{Name: "OVO"},
		{Name: "DANA"},
		{Name: "ShopeePay"},
		{Name: "LinkAja"},
		{Name: "PayPal"},
		{Name: "Apple Pay"},
		{Name: "Google Pay"},
		{Name: "Samsung Pay"},
		{Name: "Virtual Account - BCA"},
		{Name: "Virtual Account - Mandiri"},
		{Name: "Virtual Account - BNI"},
		{Name: "Virtual Account - BRI"},
		{Name: "Kredivo"},
		{Name: "Akulaku"},
		{Name: "Indomaret"},
		{Name: "Alfamart"},
	}

	for _, method := range paymentMethods {
		var existing models.PaymentMethod
		if err := db.DB.Where("name = ?", method.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.DB.Create(&method).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
