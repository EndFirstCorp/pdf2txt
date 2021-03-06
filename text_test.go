package pdf2txt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/EndFirstCorp/pdflib"
	"github.com/EndFirstCorp/peekingReader"
	"github.com/ledongthuc/pdf"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	updf "github.com/unidoc/unidoc/pdf/model"
)

func TestText(t *testing.T) {
	f, _ := os.Open(`testData/Kicker.pdf`)

	_, err := Text(f)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(r.(*bytes.Buffer).String())
}

func TestSamsung(t *testing.T) {
	f, _ := os.Open(`testData/samsung.pdf`)

	_, err := Text(f)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(r.(*bytes.Buffer).String())
}

func TestFindCatalog(t *testing.T) {
	f, _ := os.Open(`testData/102725-PUB-Replacement-PUBLIC.pdf`)

	_, err := Text(f)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(r.(*bytes.Buffer).String())
}

func TestGetTextsection(t *testing.T) {
	b, _ := ioutil.ReadFile(`testData/contents-v1.4.txt`)

	s, err := getTextSections(peekingReader.NewMemReader(b))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}

func TestGetText(t *testing.T) {
	b, _ := ioutil.ReadFile(`testData/132_0.txt`)

	s, err := getTextSections(peekingReader.NewMemReader(b))
	if err != nil {
		t.Fatal(err)
	}
	if s != nil {

	}
}

func TestProfoto(t *testing.T) {
	f, _ := os.Open(`testData/Profoto.pdf`)

	_, err := Text(f)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(r.(*bytes.Buffer).String())
}

func TestGetTextSections(t *testing.T) {
	b, _ := ioutil.ReadFile(`testData/textSection.txt`)
	_, err := getTextSections(peekingReader.NewMemReader(b))
	if err != nil {
		t.Fatal(err)
	}
}

func TestProfotoUG(t *testing.T) {
	f, _ := os.Open(`testData/ProfotoUserGuide.pdf`)

	_, err := Text(f)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(r.(*bytes.Buffer).String())
}

func TestUnexpectedEOF(t *testing.T) {
	f, _ := os.Open(`testData/financial_accounting.pdf`)

	_, err := Text(f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCh5thru8(t *testing.T) {
	f, _ := os.Open(`testData/ch5thru8.pdf`)
	_, err := Text(f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetObjectStream(t *testing.T) {
	b, _ := ioutil.ReadFile(`testData/objectstream.txt`)
	o := &object{dict: dictionary{"/Type": name("/ObjStm"), "/N": token("5"), "/First": token("34")}}
	o.stream = b
	_, err := o.getObjectStream()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCmap(t *testing.T) {
	f, _ := os.Open(`testData/bfrange.txt`)
	r := peekingReader.NewBufReader(f)
	_, err := getCmap(r)
	if err != nil {
		t.Error("expected success", err)
	}
}

func BenchmarkUnidoc(t *testing.B) {
	f, err := os.Open(`testData/Kicker.pdf`)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	pdfReader, err := updf.NewPdfReader(f)
	if err != nil {
		t.Fatal(err)
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		t.Fatal(err)
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			t.Fatal(err)
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			t.Fatal(err)
		}

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			t.Fatal(err)
		}

		// If the value is an array, the effect shall be as if all of the streams in the array were concatenated,
		// in order, to form a single stream.
		pageContentStr := ""
		for _, cstream := range contentStreams {
			pageContentStr += cstream
		}

		cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)
		_, err = cstreamParser.ExtractText()
		if err != nil {
			t.Fatal(err)
		}

		//fmt.Printf("%s\n", txt)
	}
}

func BenchmarkRscPdf(t *testing.B) {
	pdf, _ := pdf.Open(`testData/Kicker.pdf`)
	pdf.GetPlainText()
	//fmt.Println(r.(*bytes.Buffer).String())
}

func BenchmarkEndFirst(t *testing.B) {
	f, _ := os.Open(`testData/Kicker.pdf`)
	Text(f)
}

func BenchmarkPdfLib(t *testing.B) {
	r, err := pdflib.ExtractText(`testData/Kicker.pdf`, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r.(*bytes.Buffer).String())
	/*	ctx, _ := pdflib.Read(`testData/ProfotoUserGuide.pdf`, nil)
		r, err := extract.Text(ctx)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(r.(*bytes.Buffer).String())*/
}
