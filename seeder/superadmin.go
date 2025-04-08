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
