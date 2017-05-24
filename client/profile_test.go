package client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestParseProfile(t *testing.T) {
	var profile Profile
	profileFile, err := os.Open("test-profile.json")
	assert.NoError(t, err)

	err = json.NewDecoder(profileFile).Decode(&profile)
	assert.NoError(t, err)

	expectedCommander := Commander{
		Name:          "FooBar",
		Id:            290,
		Credits:       2749997,
		Debt:          0,
		Alive:         true,
		Docked:        true,
		CurrentShipId: 1,
		Rank: Rank{
			Service:    0,
			Crime:      1,
			Explore:    3,
			Federation: 0,
			Empire:     0,
			Combat:     2,
			Power:      0,
			Trade:      4,
			CQC:        0,
		},
	}
	assert.Equal(t, expectedCommander, profile.Commander)

	expectedShips := map[StringInt]Ship{
		0: Ship{
			Name: "SideWinder",
			Id:   0,
			Free: true,
			Station: Station{
				Name: "Cleve Hub",
				Id:   3230448384,
			},
			Value: ShipValue{
				Cargo:    0,
				Modules:  5934,
				Unloaned: 0,
				Hull:     0,
				Total:    5934,
			},
			Starsystem: Starsystem{
				SystemAddress: 5856221467362,
				Name:          "Eravate",
				Id:            5856221467362,
			},
		},
		1: Ship{
			Name: "CobraMkIII",
			Id:   1,
			Free: false,
			Station: Station{
				Name: "Russell Ring",
				Id:   3230448640,
			},
			Value: ShipValue{
				Cargo:    0,
				Modules:  1829408,
				Unloaned: 94845,
				Hull:     205287,
				Total:    2034695,
			},
			Starsystem: Starsystem{
				SystemAddress: 5856221467362,
				Name:          "Eravate",
				Id:            5856221467362,
			},
		},
	}
	assert.Equal(t, expectedShips, profile.Ships)

	assert.Equal(t, 13, len(profile.LastStarport.Commodities))

	expectedCommodity := Commodity{
		Name:               "Aquaponic Systems",
		Id:                 128049230,
		SellPrice:          144,
		BuyPrice:           155,
		Capacity:           38891,
		Demand:             1,
		Stock:              26127,
		SecIllegalMin:      1.22,
		SecIllegalMax:      2.78,
		CostMin:            232.62569565217,
		CostMean:           314.00,
		CostMax:            415.00,
		CategoryName:       "Technology",
		DemandBracket:      0,
		StatusFlags:        []string{},
		TargetStock:        26066,
		BaseConsumptionQty: 0,
		ConsumptionQty:     0,
		BaseCreationQty:    32.8,
		CreationQty:        26066,
		ConsumeBuy:         4,
		HomeBuy:            56,
		HomeSell:           52,
		StolenMod:          0.75,
		RareMinStock:       0,
		RareMaxStock:       0,
		VolumeScale:        1.12,
		StockBracket:       2,
		MarketId:           "",
	}
	assert.Equal(t, expectedCommodity, profile.LastStarport.Commodities[0])
}

func TestUnmarshalStringFloat(t *testing.T) {
	type TestStringFloat struct {
		MyFloat StringFloat `json:"stringfloat"`
	}

	const floatJson string = "{ \"stringfloat\": 1.234 }"
	var parsedFloatJson TestStringFloat
	err := json.NewDecoder(strings.NewReader(floatJson)).Decode(&parsedFloatJson)
	assert.NoError(t, err)
	assert.EqualValues(t, 1.234, parsedFloatJson.MyFloat)

	const stringJson string = "{ \"stringfloat\": \"2.234\" }"
	var parsedStringJson TestStringFloat
	err = json.NewDecoder(strings.NewReader(stringJson)).Decode(&parsedStringJson)
	assert.NoError(t, err)
	assert.EqualValues(t, 2.234, parsedStringJson.MyFloat)

	const notValidFloat string = "{ \"stringfloat\": \"2.twentythree4\" }"
	var parsedNotValidFloat TestStringFloat
	err = json.NewDecoder(strings.NewReader(notValidFloat)).Decode(&parsedNotValidFloat)
	assert.Error(t, err)
}

func TestUnmarshalStringInt(t *testing.T) {
	type TestStringInt struct {
		MyInt StringInt `json:"stringint"`
	}

	const intJson string = "{ \"stringint\": 14 }"
	var parsedIntJson TestStringInt
	err := json.NewDecoder(strings.NewReader(intJson)).Decode(&parsedIntJson)
	assert.NoError(t, err)
	assert.EqualValues(t, 14, parsedIntJson.MyInt)

	const stringJson string = "{ \"stringint\": \"224\" }"
	var parsedStringJson TestStringInt
	err = json.NewDecoder(strings.NewReader(stringJson)).Decode(&parsedStringJson)
	assert.NoError(t, err)
	assert.EqualValues(t, 224, parsedStringJson.MyInt)

	const notValidInt string = "{ \"stringint\": \"twentythree\" }"
	var parsedNotValidInt TestStringInt
	err = json.NewDecoder(strings.NewReader(notValidInt)).Decode(&parsedNotValidInt)
	assert.Error(t, err)
}
