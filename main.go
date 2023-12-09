package main

import (
	"fmt"
	"github.com/dslipak/pdf"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"io/ioutil"
	"strings"
)

// Example of an offline perpetual license key.
const offlineLicenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
a45d25944d22c67e4f68343991fd42adlcd7ea3845bbe3601f3c219a440clac9
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// The customer name needs to match the entry that is embedded in the signed key.
	customerName := `My Company`

	// Good to load the license key in `init`. Needs to be done prior to using the library, otherwise operations
	// will result in an error.
	err := license.SetLicenseKey(offlineLicenseKey, customerName)
	if err != nil {
		fmt.Println(err)
	}
}

func readPdf(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}

	var formattedText strings.Builder

	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		texts := p.Content().Text
		for _, text := range texts {
			if isSameSentence(text, lastTextStyle) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				formattedText.WriteString(fmt.Sprintf("%s\n", lastTextStyle.S))
				lastTextStyle = text
			}
		}
	}

	return formattedText.String(), nil
}

func readDocx(path string) (string, error) {
	doc, err := document.Open(path)
	if err != nil {
		return "", err
	}

	var formattedText strings.Builder

	for _, para := range doc.Paragraphs() {
		formattedText.WriteString(fmt.Sprintf("%s\n", para))
	}

	return formattedText.String(), nil
}

func readTxt(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func isSameSentence(text1, text2 pdf.Text) bool {
	return text1.Font == text2.Font &&
		text1.FontSize == text2.FontSize &&
		text1.X == text2.X &&
		text1.Y == text2.Y
}

func main() {

	lk := license.GetLicenseKey()
	if lk == nil {
		fmt.Printf("Failed retrieving license key")
		return
	}
	fmt.Printf("License: %s\n", lk.ToString())

	content, err := readPdf("try.pdf")
	if err != nil {
		panic(err)
	}
	fmt.Println("PDF Content:\n", content)

	content, err = readDocx("try.docx")
	if err != nil {
		panic(err)
	}
	fmt.Println("\nDOCX Content:\n", content)

	content, err = readTxt("example.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("\nTXT Content:\n", content)
}
