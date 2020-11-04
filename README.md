# imgproxyurl

imgproxyurl is a small library to help you build urls for [imgproxy](https://github.com/imgproxy/imgproxy).

This is a WIP.

### Usage
Have a look at an example to understand what's it all about:
```go
url, err := imgproxyurl.
		New("/path/to/image.jpg").
		SetKey("e99bd6542067de7dac460558ecada3987dd2d18b066180eaa1c3abc66fb22e463d177ac8f64c93c44d0d78c35adcdda7e0b5f5a116b23ac3d1fa7a305d0727c4").
		SetSalt("a997d51b78d28ba8c05f39b6e634a044b9551352b105f70a4c0fc4c0eca5982719a33527d0253810273bf4d8b747a261cd4898d3e46916cc57d1de8aac132870").
		SetWidth(400).
		SetHeight(300).
		SetResizingType(imgproxyurl.ResizingTypeFit).
		GetAbsolute("http://localhost:8080")

// url = "http://localhost:8080/a3eK6TO-pMwXvXtakEZjTov3qDrUoDeGL1Xb_1p-Ue4/w:400/h:300/rt:fit/L3BhdGgvdG8vaW1hZ2UuanBn"
```

Load key and salt from environment variables `IMGPROXY_KEY` and `IMGPROXY_SALT`:
```go
url, err := imgproxyurl.
		NewFromEnvironment("/path/to/image.jpg").
		SetWidth(400).
		SetHeight(300).
		SetResizingType(imgproxyurl.ResizingTypeFit).
		Get()

// url = "/a3eK6TO-pMwXvXtakEZjTov3qDrUoDeGL1Xb_1p-Ue4/w:400/h:300/rt:fit/L3BhdGgvdG8vaW1hZ2UuanBn"
```

### Supported processing options
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
- [quality](https://docs.imgproxy.net/#/generating_the_url_advanced?id=quality)
- [max bytes](https://docs.imgproxy.net/#/generating_the_url_advanced?id=max-bytes)
- [background](https://docs.imgproxy.net/#/generating_the_url_advanced?id=background)
- [background alpha](https://docs.imgproxy.net/#/generating_the_url_advanced?id=background-alpha)

Not all options are supported at the moment.