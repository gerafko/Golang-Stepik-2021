package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	easyjson "github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

type User struct {
	Username string   `json:"name,string"`
	Email    string   `json:"email,string"`
	Browser  []string `json:"browsers,string"`
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	uniqueBrowsers := 0
	seenBrowsers := make(map[string]int)
	var chunk []byte
	var eol bool
	fmt.Fprintln(out, "found users:")
	for i := 0; fileScanner.Scan(); i++ {
		chunk = fileScanner.Bytes()
		user := User{}
		if !eol {
			err := easyjson.Unmarshal(chunk, &user)
			if err != nil {
				panic(err)
			}
		}

		okM := false
		okA := false
		isAndroid := false
		isMSIE := false

		for _, browserRaw := range user.Browser {
			okM = strings.Contains(browserRaw, "MSIE")
			okA = strings.Contains(browserRaw, "Android")
			if okM || okA {
				if okA {
					isAndroid = true
				} else {
					isMSIE = true
				}
				if _, ok := seenBrowsers[browserRaw]; !ok {
					seenBrowsers[browserRaw] = 1
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		} else {
		}

		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", i, user.Username, email))
	}
	fmt.Fprintln(out, "\nTotal unique browsers", uniqueBrowsers)
}

var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonB8df9358DecodeHw3(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Username = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browser = nil
			} else {
				in.Delim('[')
				if out.Browser == nil {
					if !in.IsDelim(']') {
						out.Browser = make([]string, 0, 4)
					} else {
						out.Browser = []string{}
					}
				} else {
					out.Browser = (out.Browser)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browser = append(out.Browser, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonB8df9358EncodeHw3(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"browsers\":"
		out.RawString(prefix)
		if in.Browser == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browser {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB8df9358EncodeHw3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB8df9358EncodeHw3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB8df9358DecodeHw3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB8df9358DecodeHw3(l, v)
}
