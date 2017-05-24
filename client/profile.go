package client

import (
	"strconv"
	"strings"
)

type Rank struct {
	Combat     int `json:"combat"`
	Trade      int `json:"trade"`
	Explore    int `json:"explore"`
	Crime      int `json:"crime"`
	Service    int `json:"service"`
	Empire     int `json:"empire"`
	Federation int `json:"federation"`
	Power      int `json:"power"`
	CQC        int `json:"cqc"`
}

type Commander struct {
	Name          string `json:"name"`
	Id            int    `json:"id"`
	Credits       int    `json:"credits"`
	Debt          int    `json:"debt"`
	CurrentShipId int    `json:"currentShipId"`
	Alive         bool   `json:"alive"`
	Docked        bool   `json:"docked"`
	Rank          Rank   `json:"rank"`
}

type Starsystem struct {
	SystemAddress int    `json:"systemaddress"`
	Name          string `json:"name"`
	Id            int    `json:"id"`
}

type ShipValue struct {
	Cargo    int `json:"cargo"`
	Modules  int `json:"modules"`
	Unloaned int `json:"unloaned"`
	Hull     int `json:"hull"`
	Total    int `json:"total"`
}

type Ship struct {
	Name    string `json:"name"`
	Id      int    `json:"id"`
	Free    bool   `json:"free"`
	Station struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	} `json:"station"`
	Value      ShipValue  `json:"value"`
	Starsystem Starsystem `json:"starsystem"`
}

type LastSystem struct {
	Name string `json:"name"`

	// This is a string that looks like it is always an integer.
	Id StringInt `json:"id"`

	Faction string `json:"faction"`
}

type Module struct {
	Name     string `json:"name"`
	Id       int    `json:"id"`
	On       bool   `json:"on"`
	Health   int    `json:"health"`
	Value    int    `json:"value"`
	Free     bool   `json:"free"`
	Priority int    `json:"priority"`
}

type CurrentShip struct {
	Name            string     `json:"name"`
	Id              int        `json:"id"`
	Starsystem      Starsystem `json:"starsystem"`
	Free            bool       `json:"free"`
	OxygenRemaining int        `json:"oxygenRemaining"`
	Health          struct {
		ShieldUp  bool `json:"shieldup"`
		Shield    int  `json:"shield"`
		Hull      int  `json:"hull"`
		Integrity int  `json:"integrity"`
		Paintwork int  `json:"paintwork"`
	}
	Value           ShipValue `json:"value"`
	Alive           bool      `json:"alive"`
	CockpitBreached bool      `json:"cockpitBreached"`
	Station         struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	} `json:"station"`
	Modules map[string]Module `json:"modules"`
}

type SaleModule struct {
	Name     string `json:"name"`
	Id       int    `json:"id"`
	Category string `json:"category"`
	Cost     int    `json:"cost"`

	//Is this always null?
	Sku string `json:"sku"`
}

type Commodity struct {
	Name      string    `json:"name"`
	Id        StringInt `json:"id"`
	SellPrice int       `json:"sellPrice"`
	BuyPrice  int       `json:"buyPrice"`
	Capacity  int       `json:"capacity"`
	Demand    int       `json:"demand"`
	Stock     int       `json:"stock"`

	SecIllegalMin      StringFloat `json:"sec_illegal_min"`
	SecIllegalMax      StringFloat `json:"sec_illegal_max"`
	CostMin            StringFloat `json:"cost_min"`
	CostMean           StringFloat `json:"cost_mean"`
	CostMax            StringFloat `json:"cost_max"`
	CategoryName       string      `json:"categoryname"`
	DemandBracket      int         `json:"demandBracket"`
	StatusFlags        []string    `json:"statusFlags"`
	TargetStock        int         `json:"targetStock"`
	BaseConsumptionQty float64     `json:"baseConsumptionQty"`
	ConsumptionQty     int         `json:"consumptionQty"`
	BaseCreationQty    float64     `json:"baseCreationQty"`
	CreationQty        int         `json:"creationQty"`
	ConsumeBuy         StringInt   `json:"consumebuy"`
	HomeBuy            StringInt   `json:"homebuy"`
	HomeSell           StringInt   `json:"homesell"`
	StolenMod          StringFloat `json:"stolenmod"`
	RareMinStock       StringInt   `json:"rare_min_stock"`
	RareMaxStock       StringInt   `json:"rare_max_stock"`
	VolumeScale        StringFloat `json:"volumescale"`
	StockBracket       int         `json:"stockBracket"`

	// Is this always null?
	MarketId string `json:"market_id"`
}

type LastStarport struct {
	Name string `json:"name"`

	// This is a string that looks like it is always an integer.
	Id StringInt `json:"id"`

	Modules map[string]SaleModule `json:"modules"`

	Faction string `json:"faction"`

	Commodities []Commodity `json:"commodities"`
}

type Profile struct {
	Commander    Commander    `json:"commander"`
	LastSystem   LastSystem   `json:"lastSystem"`
	LastStarport LastStarport `json:"lastStarport"`

	// Each entry in the ship map is labeled with an index as a string, e.g.
	// "ships" : {
	//   "5" : { "name": "Vulture", "id" : 5, ... }
	//   ...
	//  }
	//
	// FIXME(andrew): It'd be better to convert this into an array using a custom
	// JSON deserializer, so that len(ships) and similar works.
	Ships map[StringInt]Ship `json:"ships"`

	Ship CurrentShip `json:"ship"`
}

func (c *EDClient) GetProfile(out *Profile) error {
	return c.GetJSON(DefaultProfile, out)
}

type StringFloat float64

func (x *StringFloat) UnmarshalJSON(value []byte) error {
	trimmedValue := strings.Trim(string(value), "\"")
	xx, err := strconv.ParseFloat(trimmedValue, 32)
	*x = StringFloat(xx)
	return err
}

type StringInt int

func (x *StringInt) UnmarshalJSON(value []byte) error {
	trimmedValue := strings.Trim(string(value), "\"")
	xx, err := strconv.ParseInt(trimmedValue, 0, 64)
	*x = StringInt(xx)
	return err
}
