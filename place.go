package main

type Address string
type Category string
type Radius int32 // Meters
type Place string

type PlaceAPI interface {
	Categories() []Category
	Get(PlaceOptions, Category) Place
}

type PlaceOptions struct {
	Address `json:"address,omitempty"`
	Radius  `json:"radius,omitempty"`
}

type MockPlaceAPI struct{}

func (mp MockPlaceAPI) Categories() (result []Category) {
	return []Category{"BBQ", "Pizza", "Burger", "Salad", "Breakfast"}
}

func (mp MockPlaceAPI) Get(po PlaceOptions, c Category) (result Place) {
	return Place(c)
}

// TODO: Finish google place api integration
type GooglePlaceAPI struct {
	APIKey string
}

func (gp GooglePlaceAPI) Categories() (result []Category) {
	return []Category{"BBQ", "Pizza", "Burger", "Salad", "Breakfast"}
}

func (gp GooglePlaceAPI) Get(PlaceOptions, Category) (result Place) {
	return "Winner"
}
