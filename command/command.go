package command

import (
	"reflect"
	// "encoding/json"
	"fmt"
	"os/exec"
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
		cmd.mem = append(cmd.mem, "Get-WmiObject win32_OperatingSystem |%{\"TotalPhysicalMemory: {0}`nFreePhysicalMemory : {1}`nTotalVirtualMemory : {2}`nFreeVirtualMemory  : {3}\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory}")
		cmd.disk = append(cmd.disk, "Powershell.exe")
		cmd.disk = append(cmd.disk, "get-WmiObject win32_logicaldisk")

		cpuOut = execCmd(cmd.cpu, os)
		memOut = execCmd(cmd.mem, os)
		diskOut = execCmd(cmd.disk, os)

		cpuArrOut := strings.Split(string(cpuOut), "\r\n")
		memArrOut := strings.Split(string(memOut), "\n")
		diskArrOut := strings.Split(string(diskOut), "\r\n")

		cpuResult = getMapList(cpuArrOut)
		fmt.Println("cpu###", cpuResult)
		memResult = getMapList(memArrOut)
		fmt.Println("mem###", memResult)
		diskResult = getMapList(diskArrOut)
		fmt.Println("disk###", diskResult)
	case "darwin":
		fmt.Println("MAC operating system")
	case "linux":
		fmt.Println("Linux")
		cmd.cpu = append(cmd.cpu, "bash")
		cmd.cpu = append(cmd.cpu, "-c")
		cmd.cpu = append(cmd.cpu, "top -b -n1 | grep -Po '[0-9.]+ id' | awk '{print 100-$1}'")
		cmd.mem = append(cmd.mem, "free")
		cmd.mem = append(cmd.mem, "-k")
		cmd.disk = append(cmd.disk, "df")
		cmd.disk = append(cmd.disk, "-k")

		cpuOut, err := exec.Command(cmd.cpu[0],cmd.cpu[1],cmd.cpu[2]).Output()
		if err != nil {
			log.Fatal(err)
		}
		memOut = execCmd(cmd.mem, os)
		diskOut = execCmd(cmd.disk, os)

		memArrOut := strings.Split(string(memOut), "\n")
		diskArrOut := strings.Split(string(diskOut), "\n")
		//cpu
		var cpuArr []string
		cpuArr = append(cpuArr, string(cpuOut))
		cpuResult = cpuArr
		//memory
		memMatrix := getMatrix(memArrOut)
		memMap := make(map[string]string)
		for i, _ := range memMatrix[0] {
            memMap[memMatrix[0][i]] = memMatrix[1][i+1]
		}
		//disk
		diskMatrix := getMatrix(diskArrOut)
		tmpDiskMap := make(map[string]string)

		var diskArr []interface{}
        for i:=1; i<len(diskMatrix[1]); i++ {
			for j, _ := range diskMatrix[i] {
					tmpDiskMap[diskMatrix[0][j]] = diskMatrix[i][j]
					fmt.Println(">> disk ",diskMatrix[0][j]," / ",diskMatrix[i][j])
			}
		diskArr = append(diskArr, &tmpDiskMap)
        }
		cpuResult = cpuArr
		memResult = memMap
		diskResult = diskArr
	default:
		fmt.Println("This OS is not supported : ", os)
	}

	m := make(map[string]interface{}) 
	m["cpu"] = cpuResult
	m["mem"] = memResult
	m["disk"] = diskResult
	fmt.Println("#### map : ", m)
	return m
}

func getMatrix(out []string) [][]string {
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
	var tmpArr []interface{}
	for _, val := range out {
		if val == "" {
			tmpArr = append(tmpArr, m)
			// m2[strconv.Itoa(len(m2)+1)] = m
			m = make(map[string]string)
		} else {
			tmpEl := strings.Split(val, ":")
			tmp1 := strings.Trim(tmpEl[0], " ")
			tmp2 := strings.Trim(tmpEl[1], " ")
			m[tmp1] = string(tmp2)
		}
	}
	// for key, val := range m2 {
	// 	fmt.Println(key, val)
	// }
	var resultArr []interface{}
	for _, val := range tmpArr {
		if reflect.ValueOf(val).Len() != 0 {
			resultArr = append(resultArr, val)
		} 
	}
	return resultArr
}

func execCmd(cmd []string, os string) []byte {
	out, err := exec.Command(cmd[0],cmd[1]).Output()
		if err != nil {
			log.Fatal(err)
		}
	// fmt.Printf("The output is %s\n", out)
	return out
}
