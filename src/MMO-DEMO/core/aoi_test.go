package core_test

import (
	"fmt"
	"testing"
	"zinx/MMO-DEMO/core"
)

func TestAOIManager(t *testing.T) {
	aoi := core.NewAOIManager(0, 250, 5, 0, 250, 5)
	fmt.Println(aoi)
}

func TestGetSurroundGridsByGid(t *testing.T) {
	aoi := core.NewAOIManager(0, 250, 5, 0, 250, 5)
	// fmt.Println(aoi)
	fmt.Println(aoi.GetSurroundGridsByGid(12))
}
