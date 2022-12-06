package tbclient

func (c *Client) getInventoryVmService(topologyUid string) *collectionService[InventoryVm, inventoryVmCollection] {
	return &collectionService[InventoryVm, inventoryVmCollection]{
		client:       c,
		resourcePath: "/inventory-template-vms",
		topologyUid:  topologyUid,
	}
}

func (c *Client) getVmService(topologyUid string) *resourceService[Vm, vmCollection] {
	return &resourceService[Vm, vmCollection]{
		collectionService[Vm, vmCollection]{
			client:       c,
			resourcePath: "/vms",
			topologyUid:  topologyUid,
		},
	}
}

func (c *Client) GetAllInventoryVms(topologyUid string) ([]InventoryVm, error) {
	return c.getInventoryVmService(topologyUid).getAll()
}

func (c *Client) GetAllVms(topologyUid string) ([]Vm, error) {
	return c.getVmService(topologyUid).getAll()
}

func (c *Client) GetVm(uid string) (*Vm, error) {
	return c.getVmService("").getOne(uid)
}

func (c *Client) CreateVm(vm Vm) (*Vm, error) {
	return c.getVmService("").create(vm)
}

func (c *Client) UpdateVm(vm Vm) (*Vm, error) {
	return c.getVmService("").update(vm.Uid, vm)
}

func (c *Client) DeleteVm(vmUid string) error {
	return c.getVmService("").delete(vmUid)
}
