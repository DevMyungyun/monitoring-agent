package command

import (
	// "reflect"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"runtime"
	log "github.com/sirupsen/logrus"
)

type command struct {
	cpu []string
	mem []string
	disk []string
}

func DetectOS() string {
	// OS check
	os := runtime.GOOS
	return os
}

func GetResource(os string) interface{} {
	path, _ := exec.LookPath("go")
	fmt.Println(path)

	var cmd = command{}

	var cpuOut []byte
	var memOut []byte
	var diskOut []byte

	var cpuResult interface{}
	var memResult interface{}
	var diskResult interface{}
	switch os {
	case "windows":
		fmt.Println("Windows")
		cmd.cpu = append(cmd.cpu, "Powershell.exe")
		cmd.cpu = append(cmd.cpu, "Get-WmiObject Win32_Processor | Measure-Object -Property LoadPercentage -Average")
		cmd.mem = append(cmd.mem, "Powershell.exe")
		cmd.mem = append(cmd.mem, "Get-WmiObject win32_OperatingSystem |%{\"Total Physical Memory: {0}KB`nFree Physical Memory : {1}KB`nTotal Virtual Memory : {2}KB`nFree Virtual Memory  : {3}KB\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory}")
		cmd.disk = append(cmd.disk, "Powershell.exe")
		cmd.disk = append(cmd.disk, "get-WmiObject win32_logicaldisk")

		cpuOut = execCmd(cmd.cpu, os)
		memOut = execCmd(cmd.mem, os)
		diskOut = execCmd(cmd.disk, os)

		cpuArrOut := strings.Split(string(cpuOut), "\r\n")
		memArrOut := strings.Split(string(memOut), "\n")
		diskArrOut := strings.Split(string(diskOut), "\r\n")

		cpuResult = getMapList(cpuArrOut)
		memResult = getMapList(memArrOut)
		diskResult = getMapList(diskArrOut)
	case "darwin":
		fmt.Println("MAC operating system")
	case "linux":
		fmt.Println("Linux")
		cmd.cpu = append(cmd.cpu, "top")
		cmd.cpu = append(cmd.cpu, "-b -n1 | grep -Po '[0-9.]+ id' | awk '{print 100-$1}'")
		cmd.mem = append(cmd.mem, "free")
		cmd.mem = append(cmd.mem, "-h")
		cmd.disk = append(cmd.disk, "df")
		cmd.disk = append(cmd.disk, "-h")

		cpuOut = execCmd(cmd.cpu, os)
		memOut = execCmd(cmd.mem, os)
		diskOut = execCmd(cmd.disk, os)

		// cpuArrOut := strings.Split(string(cpuOut), "\n")
		memArrOut := strings.Split(string(memOut), "\n")
		diskArrOut := strings.Split(string(diskOut), "\n")

		memResult = generateMatrix(memArrOut)
		diskResult = generateMatrix(diskArrOut)

		// cpuResult =
		memResult = generateMatrix(memArrOut)
		diskResult = generateMatrix(diskArrOut)
	default:
		fmt.Println("This OS is not supported : ", os)
	}

	m := make(map[string]interface{}) 
	m["cpu"] = cpuResult
	m["mem"] = memResult
	m["disk"] = diskResult

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	// JSON 바이트를 문자열로 변경
	jsonString := string(jsonBytes)


	// matrix := generateMatrix(arrOut)
	// for _, row := range matrix {
	// 	for i, element := range row {
	// 		if i == 3 {
	// 			fmt.Print(element, " ")
	// 		}
	// 	}
	// 	fmt.Println()
	// }
	return jsonString
}

func generateMatrix(out []string) [][]string {
	matrix := make([][]string, len(out))
	for i := 0; i < len(out); i++ {
		// fmt.Println("tmp>>>", tmpArr)
		matrix[i] = make([]string, 0, len(out))
		// vector := make([]string, len(out))
		tmpArr := strings.Split(out[i], " ")
		for j := 0; j < len(tmpArr); j++ {
			if tmpArr[j] != "" {
				matrix[i] = append(matrix[i], tmpArr[j])
				// fmt.Println("tmp>>>", tmpArr[j])
			}
		}
	}
	return matrix
}

func getMapList(out []string) interface{} {
	m := make(map[string]string)

	m2 := make(map[string]interface{})
	for _, el := range out {
		if el == "" {
			m2[strconv.Itoa(len(m2)+1)] = m
			m = make(map[string]string)
		} else {
			tmpEl := strings.Split(el, ":")
			tmp1 := strings.Trim(tmpEl[0], " ")
			tmp2 := strings.Trim(tmpEl[1], " ")
			m[tmp1] = string(tmp2)
		}
	}
	// for key, val := range m2 {
	// 	fmt.Println(key, val)
	// }
	return m2
}

func execCmd(cmd []string, os string) []byte {
	out, err := exec.Command(cmd[0],cmd[1]).Output()
		if err != nil {
			log.Fatal(err)
		}
	// fmt.Printf("The output is %s\n", out)
	return out
}
