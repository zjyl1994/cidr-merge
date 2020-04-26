package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

var convertBinTable map[string]string
var convertBinRevTable map[string]string
var convertDecTable map[string]int
var convertDecRevTable map[int]string

func main() {
	initConvTab()
	l1, err := ioutil.ReadFile("china.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l2, err := ioutil.ReadFile("china_ip_list.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	datalist := strings.Split(string(l1)+"\n"+string(l2), "\n")
	lDatalist := len(datalist)
	binDatalist := make([]string, lDatalist)
	for i := 0; i < lDatalist; i++ {
		binDatalist[i] = cidr2bin(datalist[i])
	}
	sort.Strings(binDatalist)
	combinebin(binDatalist)
	for i := 0; i < lDatalist; i++ {
		datalist[i] = bin2cidr(binDatalist[i])
	}
	writeResult(datalist)
}

func initConvTab() {
	convertBinTable = make(map[string]string, 256)
	convertBinRevTable = make(map[string]string, 256)
	convertDecTable = make(map[string]int, 256)
	convertDecRevTable = make(map[int]string, 256)
	for i := 0; i < 256; i++ {
		s := strconv.Itoa(i)
		convertBinTable[s] = fmt.Sprintf("%08b", i)
		convertBinRevTable[convertBinTable[s]] = s
		convertDecTable[s] = i
		convertDecRevTable[i] = s
	}
}

func cidr2bin(cidr string) string {
	if cidr == "" {
		return ""
	}

	scidr := strings.Split(cidr, "/") // 0 ip 1 block size
	ip := strings.Split(scidr[0], ".")
	binip := ""
	for _, v := range ip {
		binip += convertBinTable[v]
	}
	return binip[:convertDecTable[scidr[1]]]
}

func bin2cidr(bin string) string {
	if bin == "" {
		return ""
	}

	lBin := len(bin)
	bin += strings.Repeat("0", 32-lBin)
	return convertBinRevTable[bin[0:8]] + "." + convertBinRevTable[bin[8:16]] + "." + convertBinRevTable[bin[16:24]] + "." + convertBinRevTable[bin[24:32]] + "/" + convertDecRevTable[lBin]
}

func writeResult(result []string) error {
	file, err := os.Create("result.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, line := range result {
		if line != "" {
			fmt.Fprintln(w, line)
		}
	}
	return w.Flush()
}

func combinebin(bins []string) {
	lbin := len(bins)
	baseCur := 0
	for {
		if baseCur >= lbin {
			break
		}
		currentBin := bins[baseCur]
		if currentBin == "" {
			baseCur++
			continue
		}
		workCur := baseCur + 1
		for {
			if workCur >= lbin {
				break
			}
			if strings.HasPrefix(bins[workCur], currentBin) {
				bins[workCur] = ""
				workCur++
			} else {
				break
			}
		}
		baseCur = workCur
	}
}
