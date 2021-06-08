# imgproxyurl

imgproxyurl is a small library to help you build urls for [imgproxy](https://github.com/imgproxy/imgproxy).

This is a WIP.

### Usage
You can set options globally or per-instance. Create a new instance either by calling `imgproxyurl.New` or derive it from an existing one: `u.WithOptions()`.

Many imgproxy options are supported (see "Supported processing options" section below). Nevertheless, there's also a way to set an arbitrary option manually: `imgproxy.Raw{}` 
```go
// Some settings can be set globally (you'll be able to override them for specific instances)

// SetKeySalt can return an error in case key/salt is not in hex format
if err := imgproxyurl.SetKeySalt(
	"e99bd6542067de7dac460558ecada3987dd2d18b066180eaa1c3abc66fb22e463d177ac8f64c93c44d0d78c35adcdda7e0b5f5a116b23ac3d1fa7a305d0727c4",
	"a997d51b78d28ba8c05f39b6e634a044b9551352b105f70a4c0fc4c0eca5982719a33527d0253810273bf4d8b747a261cd4898d3e46916cc57d1de8aac132870",
); err != nil {
    log.Fatalln(err)
}
imgproxyurl.SetEndpoint("https://example.com/")


// create a url
u, err := imgproxyurl.New(
    "local:///o/t/otRO1jl3IUVa.jpg",
    imgproxyurl.Width{200},
    imgproxyurl.Height{200},
)
if err != nil {
    log.Fatalln(err)
}
fmt.Println(u) // https://example.com/vBTOFF_QqWqQPVCdQdjiTac8sn7EEVIh3c1UidkcvAM/h:200/w:200/bG9jYWw6Ly8vby90L290Uk8xamwzSVVWYS5qcGc


// create a copy with some options changed
u2, err := u.WithOptions(
    imgproxyurl.KeyRaw{nil},
    imgproxyurl.SaltRaw{nil},
    imgproxyurl.Format{"png"},
    imgproxyurl.PlainSourceUrl{true},
    imgproxyurl.ResizingType{imgproxyurl.ResizingTypeFill},
    imgproxyurl.Raw{OptionKey: "raw", Parameters: []interface{}{1, 2, "test"}},
)
if err != nil {
    log.Fatalln(err)
}
fmt.Println(u2) // https://example.com/insecure/h:200/raw:1:2:test/rt:fill/w:200/plain/local%3A%2F%2F%2Fo%2Ft%2FotRO1jl3IUVa.jpg@png


// playing with gravity
u3, err := u2.WithOptions(
    imgproxyurl.Gravity{
        Type:    imgproxyurl.GravityTypeFocusPoint,
        Offsets: imgproxyurl.GravityFloatOffsets{
            X: 0.3,
            Y: 0.4,
        },
    },
    imgproxyurl.Extend{
        Extend: true,
        Gravity: &imgproxyurl.Gravity{
            Type: imgproxyurl.GravityTypeNorth,
            Offsets: imgproxyurl.GravityIntegerOffsets{
                X: 100,
                Y: 200,
            },
        },
    },
)
if err != nil {
    log.Fatalln(err)
}
fmt.Println(u3) // https://example.com/insecure/ex:true:no:100:200/g:fp:0.3:0.4/h:200/rt:fill/w:200/plain/local%3A%2F%2F%2Fo%2Ft%2FotRO1jl3IUVa.jpg@png
```

### Supported processing options
You can find implementations of these processing options in `options.go`

- [resizing type](https://docs.imgproxy.net/#/generating_the_url_advanced?id=resizing-type)
- [resizing algorithm](https://docs.imgproxy.net/#/generating_the_url_advanced?id=resizing-algorithm)
- [width](https://docs.imgproxy.net/#/generating_the_url_advanced?id=width)
- [height](https://docs.imgproxy.net/#/generating_the_url_advanced?id=height)
- [dpr](https://docs.imgproxy.net/#/generating_the_url_advanced?id=dpr)
- [enlarge](https://docs.imgproxy.net/#/generating_the_url_advanced?id=enlarge)
- [extend](https://docs.imgproxy.net/#/generating_the_url_advanced?id=extend)
- [gravity](https://docs.imgproxy.net/#/generating_the_url_advanced?id=gravity)
- [crop](https://docs.imgproxy.net/#/generating_the_url_advanced?id=crop)
- [padding](https://docs.imgproxy.net/#/generating_the_url_advanced?id=padding)
- [trim](https://docs.imgproxy.net/#/generating_the_url_advanced?id=trim)
- [rotate](https://docs.imgproxy.net/#/generating_the_url_advanced?id=rotate)
- [quality](https://docs.imgproxy.net/#/generating_the_url_advanced?id=quality)
- [max bytes](https://docs.imgproxy.net/#/generating_the_url_advanced?id=max-bytes)
- [background](https://docs.imgproxy.net/#/generating_the_url_advanced?id=background)
- [background alpha](https://docs.imgproxy.net/#/generating_the_url_advanced?id=background-alpha)
- [blur](https://docs.imgproxy.net/#/generating_the_url_advanced?id=blur)
- [sharpen](https://docs.imgproxy.net/#/generating_the_url_advanced?id=sharpen)
- [preset](https://docs.imgproxy.net/#/generating_the_url_advanced?id=preset)
- [auto_rotate](https://docs.imgproxy.net/#/generating_the_url_advanced?id=auto-rotate)
- [filename](https://docs.imgproxy.net/#/generating_the_url_advanced?id=filename)

Not all options are supported at the moment.