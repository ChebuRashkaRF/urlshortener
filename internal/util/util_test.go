package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortID(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantLenStr int
	}{
		{
			name:       "simple test #1",
			url:        "https://example.com",
			wantLenStr: 8,
		},
		{
			name:       "simple test #2",
			url:        "https://yandex.ru",
			wantLenStr: 8,
		},
		{
			name:       "Empty url",
			url:        "",
			wantLenStr: 8,
		},
		{
			name: "Long url",
			url: "Hotels.AndCastles.AndHouseboats.AndIgloos.AndTeepees.AndRiversideCabins.AndLakesideCa" +
				"esOfWaterWhatsoever.AndLakeHouses.AndRegularHousesAndLodgesAndSkiLodgesAndAllThings.Ski/Cha" +
				"BungalowsAndOtherKindaLessExcitingBungalowsAndCabanasAndOceansideCabanasAndSeaSideCabanasWh",
			wantLenStr: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := GenerateShortID(tt.url)
			assert.Len(t, id, tt.wantLenStr)
		})
	}
}
