package types

type RoomStore struct {
	Rooms map[int]*Room
}

type Room struct {
	RoomIndex int
}
