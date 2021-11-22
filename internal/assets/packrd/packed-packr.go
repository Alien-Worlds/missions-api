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
	const gk = "c7892931cd55a2019eaa3a5df648ac37"
	g := packr.New(gk, "")
	hgr, err := resolver.NewHexGzip(map[string]string{
		"270824ddb8b0fdc09aa196a30242b121": "1f8b08000000000000ff74cfcf6a843010c7f1fb3cc5efa8b43e81a75843094da34d13a82789188ad47fc440f7f1177641dc5d766ec37c0ef3cd32bc4cc36f70d1c3ae446f9a33c361582139fc691d97e0032504605fdba1c765beb9164cbede1e5ddf07bf6d30fcc700509581b252c22af165f915c725bab1dda2fbf36d1c27a010ef42991d3faa6eee9ea85a8b4fa61b7cf006c9e1c594d29ce8d8572eff3351a9abfaae2f3f070000ffffc6cbbd0706010000",
		"5c7c2f25a99f8290e5426f4aa996bfff": "1f8b08000000000000ff8cd2514fc2301000e0f7fe8a7b64511234e20b4f03165dc481cd66c25353db0a8dd036ed91857f6fe6181313abf776c9d7bb5cef8643b8daeb8de7a8a07284cc6896961994e97491c15e87a0ad210302005dc6b46c3298e60f79514255e42f5576fd25a40ac26b87da9a267d4de9ec31a583dbf1386981e17b057d74e066344aa0589650548b452bdf2cf7529b0d43dd3c39356be352eef8c188edc945a532f2cce2521e3cefa688cbee53f0e8545c7a55732fcfb34764705ca8b0d58e091bf03fdd9dad958fd6448b7cc79aa2e18feee61d99b0063d17d82ce2888ac3af12ed8732ece075bfcdfbbbe4a75cd1fc39a56b78cad630e8ef2821c98490ef2738b7b521644e97abcb139c90cf000000ffffc7f374f8a9020000",
		"a263d7a869315049a59e6b6f07139d36": "1f8b08000000000000ff84d0cf4ac34010c7f1fb3cc5efd8a279829e52bb4830a6356cc09e968d1974b5990d3beb9fc71715a59243e6387cf81dbe45818b313c269f19dd4474d59ad21ad8725b1b3cc7208edf5832ad080072f2a2fe2187282e0c5f1f58736fd135d55d672ebf4d51803fa6534c9c9c1f86c4aabf0c68f6164d57d77f740caa210ace6e5b5d578d9d53791d7b4e4e9fc2a40b54b37f61974fe3f2ea0feda55fa087b6ba2ddb236ecc11abff1dd604d07a43741e7317df8568d7ee0fb3981bfa0c0000ffff129f1b7f76010000",
		"d1f3d254407d334f1ed61333f4919aef": "1f8b08000000000000ff8c924b6e83301086f7738a5906b53941563c26112a35917116595920dcc60a60845dd1e3575541316d1ef5d2f3e9f7e79959aff1a9d5ef43e9141e7a809853280845186584eab36fcca006d96a6bb5e9600588f8e75aea1a0be269983d2feae89d28dda54c20b25c203b64133905e06372d4ee540fe5e8b391318d2abb5f64f7d1566a90f6a47b7b3fd3195736d2baf2aca46bda7f925557dd21f73c7d0df9115fe888ab2b8d0a7eb0386785e0e177c0db59ce1cccbedb9c53ba63cb9400fc3671da1227165371e9f7e5415d2fe99c61421909c2827cdda5c73ce6ab1a53f1a6c53ccb95f7d9070e106c00fc0d4cccd801243cdfdfd8c0cd57000000ffff6619491fb0020000",
	})
	if err != nil {
		panic(err)
	}
	g.DefaultResolver = hgr

	func() {
		b := packr.New("migrations", "./migrations")
		b.SetResolver("001_mission.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "5c7c2f25a99f8290e5426f4aa996bfff"})
		b.SetResolver("002_explorer.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "270824ddb8b0fdc09aa196a30242b121"})
		b.SetResolver("003_explorer_mission.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "d1f3d254407d334f1ed61333f4919aef"})
		b.SetResolver("004_joinevent.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "a263d7a869315049a59e6b6f07139d36"})
	}()
	return nil
}()
