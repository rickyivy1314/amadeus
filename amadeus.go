package amadeus

import "errors"

type Amadeus struct {
	clientID     string
	clientSecret string
}

func NewAmadeusProvider(clientID, clientSecret string) *Amadeus {
	return &Amadeus{clientID: clientID, clientSecret: clientSecret}
}

func (p *Amadeus) GetFlightFees(from, to string) (map[string]interface{}, error) {

	if from == "" || to == "" {
		return nil, errors.New("invalid input")
	}
	return map[string]interface{}{
		"provider": "Amadeus",
		"fee":      100.0,
	}, nil
}
