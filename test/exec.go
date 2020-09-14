package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	path, _ := exec.LookPath("go")
	fmt.Println(path)

	out, _ := exec.Command("date").Output()
	fmt.Println(string(out))

	cmd := exec.Command("sleep", "1")
	stdoutStderr, _ := cmd.CombinedOutput()
	cmd.Start()
	cmd.Wait()
	fmt.Println(string(stdoutStderr))

	out, err := exec.Command("Powershell.exe", "get-psdrive").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The output is %s\n", out)
}
