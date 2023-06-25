package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alexmullins/zip"
	"github.com/go-pdf/fpdf"
	"github.com/go-pdf/fpdf/contrib/gofpdi"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	realgofpdi "github.com/phpdave11/gofpdi"
)

// Reference:
// - Zip with password: https://pkg.go.dev/github.com/alexmullins/zip#section-readme
// - PDF with password:
//   - https://pkg.go.dev/github.com/signintech/gopdf#section-readme
//   - PDF CPU:
//     go  get github.com/pdfcpu/pdfcpu//
//     go get github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model@v0.4.1
//     go get github.com/pdfcpu/pdfcpu/pkg/pdfcpu@v0.4.1
//     go get github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form@v0.4.1
func main() {
	protectZip()
	protectPDF2()
	fmt.Println("Done.")
}

func protectPDF() {
	inputPath := "file1.pdf"
	outputPath := "file1_protected.pdf"

	pdf := fpdf.New("P", "mm", "A4", "")
	imp := realgofpdi.NewImporter()
	imp.SetSourceFile(inputPath)

	pageCount := imp.GetNumPages()
	pageSizes := imp.GetPageSizes()

	impWrapper := gofpdi.NewImporter()

	pdf.SetAutoPageBreak(false, 0)
	pdf.SetProtection(fpdf.CnProtectPrint|fpdf.CnProtectCopy|fpdf.CnProtectModify|fpdf.CnProtectAnnotForms, "pass", "pass")

	for i := 1; i <= pageCount; i++ {
		var sizeType fpdf.SizeType
		sizeType.Wd = pageSizes[i]["/MediaBox"]["w"]
		sizeType.Ht = pageSizes[i]["/MediaBox"]["h"]

		pdf.AddPageFormat("P", sizeType)

		tplid := impWrapper.ImportPage(pdf, inputPath, i, "/MediaBox")
		impWrapper.UseImportedTemplate(pdf, tplid, 0, 0, pageSizes[i]["/MediaBox"]["w"], pageSizes[i]["/MediaBox"]["h"])

	}

	pdf.OutputFileAndClose(outputPath)

}

func protectPDF2() {

	config := model.NewAESConfiguration("pass", "pass", 256)
	config.ValidationMode = model.ValidationNone
	err := pdfcpu.EncryptFile("file1.pdf", "file1_protected.pdf", config)
	if err != nil {
		panic(err)
	}
}

func protectZip() {
	outputPath := "output.zip"
	archive, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}

	password := "pass"
	zipw := zip.NewWriter(archive)
	inputFiles := []string{
		"file1.pdf",
		"file2.doc",
	}
	for _, filePath := range inputFiles {
		w, err := zipw.Encrypt(filePath, password)
		if err != nil {
			panic(err)
		}

		f, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			panic(err)
		}
	}

	err = zipw.Close()
	if err != nil {
		panic(err)
	}

}
