package imgproxyurl

import (
	"reflect"
	"testing"
)

const testKey = "e99bd6542067de7dac460558ecada3987dd2d18b066180eaa1c3abc66fb22e463d177ac8f64c93c44d0d78c35adcdda7e0b5f5a116b23ac3d1fa7a305d0727c4"
const testSalt = "a997d51b78d28ba8c05f39b6e634a044b9551352b105f70a4c0fc4c0eca5982719a33527d0253810273bf4d8b747a261cd4898d3e46916cc57d1de8aac132870"

func TestUrl_applyOptions(t *testing.T) {
	type fields struct {
		key            []byte
		salt           []byte
		options        map[string]string
		sourceUrl      string
		plainSourceUrl bool
		format         string
		endpoint       string
		signatureSize  int
	}
	type args struct {
		options []Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Url
		wantErr bool
	}{
		{name: "key/salt", fields: fields{}, args: args{options: []Option{Key{testKey}, Salt{testSalt}}}, want: &Url{
			key:  []byte{233, 155, 214, 84, 32, 103, 222, 125, 172, 70, 5, 88, 236, 173, 163, 152, 125, 210, 209, 139, 6, 97, 128, 234, 161, 195, 171, 198, 111, 178, 46, 70, 61, 23, 122, 200, 246, 76, 147, 196, 77, 13, 120, 195, 90, 220, 221, 167, 224, 181, 245, 161, 22, 178, 58, 195, 209, 250, 122, 48, 93, 7, 39, 196},
			salt: []byte{169, 151, 213, 27, 120, 210, 139, 168, 192, 95, 57, 182, 230, 52, 160, 68, 185, 85, 19, 82, 177, 5, 247, 10, 76, 15, 196, 192, 236, 165, 152, 39, 25, 163, 53, 39, 208, 37, 56, 16, 39, 59, 244, 216, 183, 71, 162, 97, 205, 72, 152, 211, 228, 105, 22, 204, 87, 209, 222, 138, 172, 19, 40, 112},
		}, wantErr: false},
		{name: "malformed key/salt", fields: fields{}, args: args{options: []Option{Key{"key"}, Salt{"salt"}}}, want: &Url{}, wantErr: true},
		{name: "options overriding", fields: fields{
			options: map[string]string{
				"z": "50",
				"h": "100",
			},
		}, args: args{options: []Option{Width{200}, Height{300}}}, want: &Url{
			options: map[string]string{
				"z": "50",
				"h": "300",
				"w": "200",
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Url{
				key:            tt.fields.key,
				salt:           tt.fields.salt,
				options:        tt.fields.options,
				sourceUrl:      tt.fields.sourceUrl,
				plainSourceUrl: tt.fields.plainSourceUrl,
				format:         tt.fields.format,
				endpoint:       tt.fields.endpoint,
				signatureSize:  tt.fields.signatureSize,
			}
			if err := u.applyOptions(tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("applyOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(u, tt.want) {
				t.Errorf("applyOptions() error: want %v, got %v", tt.want, u)
			}
		})
	}
}

func TestUrl_clone(t *testing.T) {
	type fields struct {
		key            []byte
		salt           []byte
		options        map[string]string
		sourceUrl      string
		plainSourceUrl bool
		format         string
		endpoint       string
		signatureSize  int
	}
	type args struct {
		addOptions []Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Url
		wantErr bool
	}{
		{name: "clone w/o additional options", fields: fields{
			key:  []byte{1, 2, 3},
			salt: []byte{2, 3, 4},
			options: map[string]string{
				"a": "100",
			},
			sourceUrl:      "https://example.com/test.jpg",
			plainSourceUrl: true,
			format:         "png",
			endpoint:       "https://example2.com/test2.png",
			signatureSize:  32,
		}, args: args{}, want: &Url{
			key:  []byte{1, 2, 3},
			salt: []byte{2, 3, 4},
			options: map[string]string{
				"a": "100",
			},
			sourceUrl:      "https://example.com/test.jpg",
			plainSourceUrl: true,
			format:         "png",
			endpoint:       "https://example2.com/test2.png",
			signatureSize:  32,
		}, wantErr: false},

		{name: "clone w/ additional options", fields: fields{
			key:  []byte{1, 2, 3},
			salt: []byte{2, 3, 4},
			options: map[string]string{
				"a": "100",
			},
			sourceUrl:      "https://example.com/test.jpg",
			plainSourceUrl: true,
			format:         "png",
			endpoint:       "https://example2.com/test2.png",
			signatureSize:  32,
		}, args: args{addOptions: []Option{Width{100}, SignatureSize{31}}}, want: &Url{
			key:  []byte{1, 2, 3},
			salt: []byte{2, 3, 4},
			options: map[string]string{
				"a": "100",
				"w": "100",
			},
			sourceUrl:      "https://example.com/test.jpg",
			plainSourceUrl: true,
			format:         "png",
			endpoint:       "https://example2.com/test2.png",
			signatureSize:  31,
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Url{
				key:            tt.fields.key,
				salt:           tt.fields.salt,
				options:        tt.fields.options,
				sourceUrl:      tt.fields.sourceUrl,
				plainSourceUrl: tt.fields.plainSourceUrl,
				format:         tt.fields.format,
				endpoint:       tt.fields.endpoint,
				signatureSize:  tt.fields.signatureSize,
			}
			got, err := u.clone(tt.args.addOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("clone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clone() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrl_String(t *testing.T) {
	tests := []struct {
		name string
		u    *Url
		want string
	}{
		{name: "basic options", u: func() *Url {
			u, _ := New(
				"local:///o/t/otRO1jl3IUVa.jpg",
				Width{200},
				Height{200},
				Format{"png"},
				PlainSourceUrl{false},
				ResizingType{ResizingTypeFill},
				Key{testKey},
				Salt{testSalt},
				Endpoint{"https://example.com/"},
			)
			return u
		}(), want: "https://example.com/Yysx5pZ_gcWJbVQEHSp37U6r3swrZgFAygnHmbFK2VE/h:200/rt:fill/w:200/bG9jYWw6Ly8vby90L290Uk8xamwzSVVWYS5qcGc.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
