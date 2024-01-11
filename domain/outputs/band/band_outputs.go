package bandoutputs

import accountoutputs "github.com/mazurco066/playliter-api-go/domain/outputs/account"

type BandOutput struct {
	ID          uint                          `json:"id"`
	Logo        string                        `json:"logo"`
	Title       string                        `json:"title"`
	Description string                        `json:"description"`
	Owner       *accountoutputs.AccountOutput `json:"owner"`
}

type BandRequestOutput struct {
	ID      uint                          `json:"id"`
	Band    *BandOutput                   `json:"band"`
	Invited *accountoutputs.AccountOutput `json:"invited"`
	Status  string                        `json:"status"`
}
