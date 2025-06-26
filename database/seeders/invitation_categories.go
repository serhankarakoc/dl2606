package seeders

import (
	"davet.link/configs/logconfig"
	"davet.link/models"
	"gorm.io/gorm"
)

func SeedInvitationCategories(db *gorm.DB) error {
	categories := []models.InvitationCategory{
		{Name: "Açılış", Template: "title", Icon: "fas fa-door-open", IsActive: true},
		{Name: "After Party", Template: "title", Icon: "fas fa-glass-cheers", IsActive: true},
		{Name: "Anıtkabir Ziyareti", Template: "title", Icon: "fas fa-landmark", IsActive: true},
		{Name: "Asker Eğlencesi", Template: "person", Icon: "fas fa-user-tie", IsActive: true},
		{Name: "Baby Shower", Template: "person-family", Icon: "fas fa-baby", IsActive: true},
		{Name: "Balo", Template: "title", Icon: "fas fa-mask", IsActive: true},
		{Name: "Bekarlığa Veda", Template: "person", Icon: "fas fa-glass-martini-alt", IsActive: true},
		{Name: "Cinsiyet Partisi", Template: "person-family", Icon: "fas fa-ribbon", IsActive: true},
		{Name: "Defile", Template: "title", Icon: "fas fa-crown", IsActive: true},
		{Name: "Dini Tören", Template: "title", Icon: "fas fa-praying-hands", IsActive: true},
		{Name: "Doğum Günü", Template: "person", Icon: "fas fa-birthday-cake", IsActive: true},
		{Name: "Düğün", Template: "wedding", Icon: "fas fa-ring", IsActive: true},
		{Name: "Eğitim", Template: "title", Icon: "fas fa-book-open", IsActive: true},
		{Name: "Film Galası", Template: "title", Icon: "fas fa-film", IsActive: true},
		{Name: "Fuar", Template: "title", Icon: "fas fa-building", IsActive: true},
		{Name: "Gelin Hamamı", Template: "person", Icon: "fas fa-spa", IsActive: true},
		{Name: "Gezi", Template: "title", Icon: "fas fa-suitcase-rolling", IsActive: true},
		{Name: "Kına Gecesi", Template: "wedding", Icon: "fas fa-hand-sparkles", IsActive: true},
		{Name: "Konferans", Template: "title", Icon: "fas fa-microphone-alt", IsActive: true},
		{Name: "Kongre", Template: "title", Icon: "fas fa-users", IsActive: true},
		{Name: "Konser", Template: "title", Icon: "fas fa-music", IsActive: true},
		{Name: "Lansman", Template: "title", Icon: "fas fa-laptop", IsActive: true},
		{Name: "Mezuniyet", Template: "title", Icon: "fas fa-graduation-cap", IsActive: true},
		{Name: "Nikâh Töreni", Template: "wedding", Icon: "fas fa-ring", IsActive: true},
		{Name: "Nişan", Template: "wedding", Icon: "fas fa-heart", IsActive: true},
		{Name: "Online Etkinlik", Template: "online", Icon: "fas fa-link", IsActive: true},
		{Name: "Seminer", Template: "title", Icon: "fas fa-chalkboard-teacher", IsActive: true},
		{Name: "Sergi", Template: "title", Icon: "fas fa-palette", IsActive: true},
		{Name: "Spor Müsabakası", Template: "title", Icon: "fas fa-futbol", IsActive: true},
		{Name: "Sünnet Düğünü", Template: "person-family", Icon: "fas fa-child", IsActive: true},
		{Name: "Tanıtım", Template: "title", Icon: "fas fa-bullhorn", IsActive: true},
		{Name: "Tiyatro Gösterisi", Template: "title", Icon: "fas fa-theater-masks", IsActive: true},
		{Name: "Toplantı", Template: "title", Icon: "fas fa-calendar-day", IsActive: true},
		{Name: "Veda Partisi", Template: "title", Icon: "fas fa-calendar-times", IsActive: true},
		{Name: "Yılbaşı Partisi", Template: "title", Icon: "fas fa-calendar-alt", IsActive: true},
		{Name: "Yıldönümü", Template: "title", Icon: "fas fa-calendar-day", IsActive: true},
	}

	logconfig.SLog.Info("Davetiye kategorileri yükleniyor...")

	// Her bir kategori için
	for _, category := range categories {
		// Kategori zaten var mı kontrol et
		var existingCategory models.InvitationCategory
		if err := db.Where("name = ?", category.Name).First(&existingCategory).Error; err == gorm.ErrRecordNotFound {
			// Kategori yoksa ekle
			if err := db.Create(&category).Error; err != nil {
				logconfig.SLog.Error("Kategori eklenirken hata: " + category.Name)
				return err
			}
			logconfig.SLog.Info("Kategori eklendi: " + category.Name)
		}
	}

	logconfig.SLog.Info("Davetiye kategorileri yükleme işlemi tamamlandı.")
	return nil
}
