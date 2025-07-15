package scrapping

type Request struct {
	Url        string `json:"url"`
	PartnerTag string `json:"partner_tag"`
}
