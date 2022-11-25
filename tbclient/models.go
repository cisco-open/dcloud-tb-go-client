package tbclient

// Errors
type vndError struct {
	Logref  string `json:"logref"`
	Message string `json:"message"`
}

// Collections
type collection[T any] struct {
	Embedded T `json:"_embedded"`
}

type embeddedData[R any] interface {
	getData() []R
}

// Topology
type Topology struct {
	Uid         string `json:"uid,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Datacenter  string `json:"datacenter"`
	Notes       string `json:"notes"`
	Status      string `json:"status,omitempty"`
}

type topologyCollection struct {
	Data []Topology `json:"topologies"`
}

func (t topologyCollection) getData() []Topology {
	return t.Data
}

// InventoryNetwork
type InventoryNetwork struct {
	Id     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Subnet string `json:"subnet,omitempty"`
}

type inventoryNetworkCollection struct {
	Data []InventoryNetwork `json:"inventoryNetworks"`
}

func (i inventoryNetworkCollection) getData() []InventoryNetwork {
	return i.Data
}

// Network
type Network struct {
	Uid              string           `json:"uid,omitempty"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	InventoryNetwork InventoryNetwork `json:"inventoryNetwork"`
	Topology         Topology         `json:"topology"`
}

type networkCollection struct {
	Data []Network `json:"networks"`
}

func (n networkCollection) getData() []Network {
	return n.Data
}
