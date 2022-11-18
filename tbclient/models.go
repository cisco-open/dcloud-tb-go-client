package tbclient

type vndError struct {
	Logref  string `json:"logref"`
	Message string `json:"message"`
}

type collection[T any] struct {
	Embedded T `json:"_embedded"`
}

type embeddedData[R any] interface {
	getData() []R
}

type topologyCollection struct {
	Data []Topology `json:"topologies"`
}

func (t topologyCollection) getData() []Topology {
	return t.Data
}

type Topology struct {
	Uid         string `json:"uid,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Datacenter  string `json:"datacenter"`
	Notes       string `json:"notes"`
	Status      string `json:"status,omitempty"`
}
