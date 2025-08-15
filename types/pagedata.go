package types

type PageData struct {
	PageTitle   string
	CurrentPage string
}

type IndexPageData struct {
	PageData
	Rooms    map[int]*Room
	Username string
}

type RoomPageData struct {
	PageData
	Username       string
	Room           *Room
	MessageHistory map[string]string
}
