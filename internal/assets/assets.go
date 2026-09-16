package assets

import _ "embed"

type Asset struct {
	Name string
	Data []byte
}

//go:embed searchBarDark.png
var searchBarDark []byte

//go:embed searchBarLight.png
var searchBarLight []byte

//go:embed rewardsLogo.png
var rewardsLogoPNG []byte

//go:embed rewardsLogo.ico
var rewardsLogoICO []byte

//go:embed agreeContinue.png
var agreeContinue []byte

//go:embed readyToClaim.png
var readyToClaim []byte

//go:embed claimPoints.png
var claimPoints []byte

//go:embed closeClaimPoints.png
var closeClaimPoints []byte

//go:embed visualSearchWhite.png
var visualSearchWhite []byte

//go:embed visualSearchDark.png
var visualSearchDark []byte

//go:embed pasteImageURL.png
var pasteImageURL []byte

//go:embed ganharLabel.png
var ganharLabel []byte

//go:embed searchPlusDark.png
var searchPlusDark []byte

//go:embed searchPlusLight.png
var searchPlusLight []byte

var (
	SearchBarDark     = Asset{Name: "searchBarDark", Data: searchBarDark}
	SearchBarLight    = Asset{Name: "searchBarLight", Data: searchBarLight}
	RewardsLogoPNG    = Asset{Name: "rewardsLogoPNG", Data: rewardsLogoPNG}
	RewardsLogoICO    = Asset{Name: "rewardsLogoICO", Data: rewardsLogoICO}
	AgreeContinue     = Asset{Name: "agreeContinue", Data: agreeContinue}
	ReadyToClaim      = Asset{Name: "readyToClaim", Data: readyToClaim}
	ClaimPoints       = Asset{Name: "claimPoints", Data: claimPoints}
	CloseClaimPoints  = Asset{Name: "closeClaimPoints", Data: closeClaimPoints}
	VisualSearchWhite = Asset{Name: "visualSearchWhite", Data: visualSearchWhite}
	VisualSearchDark  = Asset{Name: "visualSearchDark", Data: visualSearchDark}
	PasteImageURL     = Asset{Name: "pasteImageURL", Data: pasteImageURL}
	GanharLabel       = Asset{Name: "ganharLabel", Data: ganharLabel}
	SearchPlusDark    = Asset{Name: "searchPlusDark", Data: searchPlusDark}
	SearchPlusLight   = Asset{Name: "searchPlusLight", Data: searchPlusLight}
)
