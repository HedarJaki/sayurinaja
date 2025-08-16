package initializer

import "mobapp/model"

func SyncDatabase() {
	DB.AutoMigrate(&model.User_app{})
}
