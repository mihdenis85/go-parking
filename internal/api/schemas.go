package api

type AddParkingSpaceLogSchema struct {
	FirstName    string `json:"first_name" binding:"required" example:"Иван"`
	LastName     string `json:"last_name" binding:"required" example:"Иванов"`
	CarMake      string `json:"car_make" binding:"required" example:"Toyota"`
	LicensePlate string `json:"license_plate" binding:"required" example:"А123БВ777"`
}
