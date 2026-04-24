package assets

import _ "embed"

type Asset struct {
	Name string
	Data []byte
}

// //go:embed searchBarDark.png
// var searchBarDark []byte

// //go:embed searchBarLight.png
// var searchBarLight []byte

//go:embed rewardsLogo.png
var rewardsLogoPNG []byte

//go:embed rewardsLogo.ico
var rewardsLogoICO []byte

//go:embed agreeContinue.png
var agreeContinue []byte

//go:embed aceitarButton.png
var aceitarButton []byte

//go:embed searchPlusDark.png
var searchPlusDark []byte

//go:embed searchPlusLight.png
var searchPlusLight []byte

//go:embed claimSearchBonus.png
var claimSearchBonus []byte

var (
	// SearchBarDark   = Asset{Name: "searchBarDark", Data: searchBarDark}
	// SearchBarLight  = Asset{Name: "searchBarLight", Data: searchBarLight}
	RewardsLogoPNG   = Asset{Name: "rewardsLogoPNG", Data: rewardsLogoPNG}
	RewardsLogoICO   = Asset{Name: "rewardsLogoICO", Data: rewardsLogoICO}
	AgreeContinue    = Asset{Name: "agreeContinue", Data: agreeContinue}
	AceitarButton    = Asset{Name: "aceitarButton", Data: aceitarButton}
	SearchPlusDark   = Asset{Name: "searchPlusDark", Data: searchPlusDark}
	SearchPlusLight  = Asset{Name: "searchPlusLight", Data: searchPlusLight}
	ClaimSearchBonus = Asset{Name: "claimSearchBonus", Data: claimSearchBonus}
)
