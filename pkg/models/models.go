package models

type Item struct {
	UUID        string
	Name        string
	Departure   string
	Destination string
	Time        string
	Capacity    int
	Passenger   int
	Passengers  []PassengerModel
}

type PassengerModel struct {
	Namelist string
	Comment  string
}
