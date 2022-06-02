package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func main() {

	templ, err := ioutil.ReadFile("./CustomerAuthorizationMaintenance.html")
	if err != nil {
		log.Fatal(err)
	}
	replaceString := strings.NewReplacer("{{CUSTOMER_NAME}}", "Tegar Abdijaya", "{{MAINTENANCE_TYPE}}", "Test", "{{EXPIRY_DATE}}", "05/06/2022")

	replaceHtml := replaceString.Replace(string(templ))
	err = ioutil.WriteFile("./CustomerAuthorizationMaintenanceChanged.html", []byte(replaceHtml), 0666)
	if err != nil {
		log.Fatal(err)
	}

	execString := []string{"./CustomerAuthorizationMaintenanceChanged.html", "./CustomerAuthorizationMaintenanceChanged.pdf"}
	execCmd := exec.Command("C:/Program Files/wkhtmltopdf/bin/wkhtmltopdf", execString...)

	execOut, err := execCmd.CombinedOutput()

	if err != nil {
		log.Fatal("Error: ", err, execOut)
	}
}
