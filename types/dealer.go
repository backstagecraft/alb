package types

// A Dealer is a service account which is reponsible for exchange of two
// different BVS assets between two independent users.
type Dealer struct {
	Id          string `json:"id"`
	Owner       string `json:"owner"`
	Replacement Asset  `json:"replacement"`
	Current     Asset  `json:"current"`
}
