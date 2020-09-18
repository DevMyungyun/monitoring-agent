package command

import (
	"fmt"
	"log"
	"os/exec"
)

func CheckDisk() {
	path, _ := exec.LookPath("go")
	fmt.Println(path)

	out, err := exec.Command("Powershell.exe", "get-psdrive").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The output is %s\n", out)
}
