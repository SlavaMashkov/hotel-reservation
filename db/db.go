package db

// TODO: get this from config
const (
	NAME     = "hotel-reservation"
	USERNAME = "root"
	PASSWORD = "example"
	URI      = "mongodb://localhost:27017"
)

type Store struct {
	UserStore  UserStore
	HotelStore HotelStore
	RoomStore  RoomStore
}
