// +build !skippackr
// Code generated by github.com/gobuffalo/packr/v2. DO NOT EDIT.

// You can use the "packr2 clean" command to clean up this,
// and any other packr generated files.
package packrd

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file/resolver"
)

var _ = func() error {
	const gk = "b156da144dcd1bf1694709c453bd270e"
	g := packr.New(gk, "")
	hgr, err := resolver.NewHexGzip(map[string]string{
		"313b57106b840836a654041eaef2dd03": "1f8b08000000000000ff8cd2cf6af3300c00f0bb9f42c786ef2b74839e7a4adbb08575696792414fc6b3bdd6acb18dad10faf6234bd33f8379d34df0b3842c8dc7f0afd63bcf5141e50859d02c2d3328d3f92a835a87a0ad2123020043c6b4ec3298e70f79514255e42f55f6ff4b481584d70eb5355dfa9ad2c5634a47f7d369d203c36b059718c0dd649240b12ea1a856ab5ebe59eea5363b86ba7b726ad6c7ad3cf0c688fdc945a532f2cce252369e0f53c4e5f02978742a2ebd6ab997e7d92332382e54d86bc7840df897eeceb6ca476ba2457e605dd1f04b77f38e4c58839e0bec167144c5e14789f64319d6787db5ee49f25d6e68fe9cd22d3c655b185dee2821c98c90eb135cdad610b2a4ebcded09ce3e030000ffff88294f84a8020000",
		"88a73843ef9e559bff744fce92ccf127": "1f8b08000000000000ff74cfcf6a843010c7f1fb3cc5efa8b43e81a75843094da34d13a82789188ad47fc440f7f1177641dc5d766ec37c0ef3cd32bc4cc36f70d1c3ae446f9a33c361582139fc691d97e0032504605fdba1c765beb9164cbede1e5ddf07bf6d30fcc700509581b252c22af165f915c725bab1dda2fbf36d1c27a010ef42991d3faa6eee9ea85a8b4fa61b7cf006c9e1c594d29ce8d8572eff3351a9abfaae2f3f070000ffffc6cbbd0706010000",
		"a6cf3d6c4ae2dbc938829a9ae4df36ea": "1f8b08000000000000ff8c924b6e83301086f7738a5906b53941563c26112a35917116595920dcc60a60845dd1e3575541316d1ef5d2f3e9f7e79959aff1a9d5ef43e9141e7a809853280845186584eab36fcca006d96a6bb5e9600588f8e75aea1a0be269983d2feae89d28dda54c20b25c203b64133905e06372d4ee540fe5e8b391318d2abb5f64f7d1566a90f6a47b7b3fd3195736d2baf2aca46bda7f925557dd21f73c7d0df9115fe888ab2b8d0a7eb0386785e0e177c0db59ce1cccbedb9c53ba63cb9400fc3671da1227165371e9f7e5415d2fe99c61421909c2827cdda5c73ce6ab1a53f1a6c53ccb95f7d9070e106c00fc0d4cccd801243cdfdfd8c0cd57000000ffff6619491fb0020000",
	})
	if err != nil {
		panic(err)
	}
	g.DefaultResolver = hgr

	func() {
		b := packr.New("migrations", "./migrations")
		b.SetResolver("001_mission.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "313b57106b840836a654041eaef2dd03"})
		b.SetResolver("002_explorer.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "88a73843ef9e559bff744fce92ccf127"})
		b.SetResolver("003_explorer_mission.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "a6cf3d6c4ae2dbc938829a9ae4df36ea"})
	}()
	return nil
}()