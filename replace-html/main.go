package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func main() {

	wkhtmltopdf.SetPath("C:/Program Files/wkhtmltopdf/bin/wkhtmltopdf") //hilangkan apabila sudah set environment di system

	fileName := "CustomerAuthorizationMaintenance"
	templ, err := ioutil.ReadFile("./" + fileName + ".html")
	if err != nil {
		log.Fatal(err)
	}
	replaceString := strings.NewReplacer(
		"{{CUSTOMER_NAME}}", "Tegar Abdijaya",
		"{{MAINTENANCE_TYPE}}", "Test",
		"{{EXPIRY_DATE}}", "05/06/2022",
		"{{EMAIL}}", "nathaniel.a@ocbcnisp.com",
		"{{SUBJECT}}", "Authorization Maintenance Test",
	)

	replaceHtml := replaceString.Replace(string(templ))
	err = ioutil.WriteFile("./"+fileName+"Changed.html", []byte(replaceHtml), 0666)
	if err != nil {
		log.Fatal(err)
	}

	// ! False Practice ## Seharusnya langsung dari []byte ke buffer lalu pdf, karena apabila melakukan proses lagi dari output html ke output pdf akan memakan resource lebih banyak
	// execString := []string{"./CustomerAuthorizationMaintenanceChanged.html", "./CustomerAuthorizationMaintenanceChanged.pdf"}
	// execCmd := exec.Command("C:/Program Files/wkhtmltopdf/bin/wkhtmltopdf", execString...)

	// execOut, err := execCmd.CombinedOutput()

	// if err != nil {
	// 	log.Fatal("Error: ", err, execOut)
	// }

	generate, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	generate.Dpi.Set(300)
	generate.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	generate.Grayscale.Set(false)

	page := wkhtmltopdf.NewPageReader(strings.NewReader(replaceHtml))
	//page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(replaceHtml)))

	page.Zoom.Set(0.95)

	generate.AddPage(page)

	errCreate := generate.Create()
	if errCreate != nil {
		log.Fatal(errCreate)
	}

	errWrite := generate.WriteFile("./" + fileName + ".pdf")
	if errWrite != nil {
		log.Fatal(errWrite)
	}

	fmt.Println("Done Create PDF")
}
