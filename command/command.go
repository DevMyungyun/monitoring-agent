package command

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func CheckDisk() []string {
	path, _ := exec.LookPath("go")
	fmt.Println(path)

	out, err := exec.Command("Powershell.exe", "get-psdrive").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The output is %s\n", out)

	stringOut := string(out)
	arrOut := strings.Split(stringOut,"\r\n")

	matrix := generateMatrix(arrOut)

	for _, row := range matrix {
		for i, element := range row {
			if(i == 3) {
				fmt.Print(element, " ")
			}
		}
		fmt.Println()
	}
	
	

	fmt.Print("end>>>", matrix)
}


func generateMatrix(out []string) [][]string {

	matrix := make([][]string, len(out))
	for i := 0; i < len(out); i++ {
		// fmt.Println("tmp>>>", tmpArr)
		matrix[i] = make([]string, 0, len(out))
		// vector := make([]string, len(out))
		tmpArr := strings.Split(out[i]," ")
		for j := 0; j<len(tmpArr); j++ {
			if tmpArr[j] != "" {
				matrix[i] = append(matrix[i], tmpArr[j])
				// fmt.Println("tmp>>>", tmpArr[j])
			} 
		}
	}

	return matrix
}