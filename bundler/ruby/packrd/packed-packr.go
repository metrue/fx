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
	const gk = "5957c0f80ccc65e6ad8f2c329312596e"
	g := packr.New(gk, "")
	hgr, err := resolver.NewHexGzip(map[string]string{
		"5d0faa4e62eb53572bfc7410d0f064f9": "1f8b08000000000000ff720bf2f755c84d2d292a4dd54fabd02d2a4daad44d4a2c4ee5e272f60f8854d053d0e3728d08f00f765530363030e072f6755100a951482c28d02b4a52d02d008b2be8e62b18e8812117200000ffff79eb4a3952000000",
		"93130b893788876c223c2a0f3de05793": "1f8b08000000000000ffccce41aac23010c6f1fd9ce2dbcda6bc577057a8579168a75a91a4ce4ca520de5ddad4780537e19ffc207c2af7695001db10836b602a2f574b91e9733fa8dc820f0f01f7f39f1e99c8c4d18c49bdc2aeae6ba2319983ff195d220038f98c16cfb58166f949ccd1eeb165f5251b5334c996bba079f0c916ca55e022a1135d65cb4caff5ece7650049ec88cef29bc3de010000ffff22ee4db77f010000",
		"973570cec900912c2cdd5c5531be70ac": "1f8b08000000000000ff3ccd410ac2301085e1fd9ce2d16c14c40308f122e2a276a674111acd4c3022de5d6262773fef8319470e931578bcc9013825796451833fa3e7e10f7a8fab4a93d69dd446cb5aa1559f171959d26fef59e1432c33e6b29bacec09f5fb653b7e3dde22bfe0312c1242c433a6c003c9caf40d0000ffff99b9e753aa000000",
	})
	if err != nil {
		panic(err)
	}
	g.DefaultResolver = hgr

	func() {
		b := packr.New("ruby", "./assets")
		b.SetResolver("Dockerfile", packr.Pointer{ForwardBox: gk, ForwardPath: "5d0faa4e62eb53572bfc7410d0f064f9"})
		b.SetResolver("app.rb", packr.Pointer{ForwardBox: gk, ForwardPath: "93130b893788876c223c2a0f3de05793"})
		b.SetResolver("fx.rb", packr.Pointer{ForwardBox: gk, ForwardPath: "973570cec900912c2cdd5c5531be70ac"})
	}()
	return nil
}()
