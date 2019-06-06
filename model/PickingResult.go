package model

type PickingResult struct {
	ItemId        int    `json:"id"`
	ItemName      string `json:"name"`
	PickedSuccess bool   `json:"pickedSuccess"`
}
