package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"bufio"
	"io"
	"strings"
	"strconv"
	"github.com/prometheus/prometheus/prompb"	
)

var promPath string
var promConfig = filepath.Join("..", "..", "documentation", "examples", "prometheus.yml")


func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	var err error
	promPath, err = os.Getwd()
	if err != nil {
		fmt.Printf("can't get current dir :%s \n", err)
		os.Exit(1)
	}
	promPath = filepath.Join(promPath, "testData/blm_prometheus.log")
	fmt.Println(promPath)
	//build := exec.Command("go", "build", "-o", promPath)
	//output, err := build.CombinedOutput()
	//if err != nil {
	//	fmt.Printf("compilation error :%s \n", output)
	//	os.Exit(1)
	//}

	exitCode := m.Run()
	//os.Remove(promPath)
	//os.RemoveAll(promData)
	os.Exit(exitCode)
}

func TestSerialization(t *testing.T) {
	var req prompb.WriteRequest;
	var ts []prompb.TimeSeries
	var tse prompb.TimeSeries
	var sample prompb.Sample

	testfile,err := os.OpenFile(promPath,os.O_RDWR,0666)
	if err != nil {
        fmt.Println("Open file error!", err)
        return
	}
	defer testfile.Close()
	fmt.Println(promPath)
	buf := bufio.NewReader(testfile)
	i :=0
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok! line:", i)
				break
			} else {
				fmt.Println("Read file error!",err)
				return 
			}
		}
		if strings.Contains(line,"server.go:201:")  {
			sa := strings.Split(line," ")
			sample.Timestamp, _ = strconv.ParseInt(sa[7][:(len(sa[7])-1)],10,64)
			sample.Value,_  = strconv.ParseFloat(sa[9][:(len(sa[9])-1)],64)
			tse.Samples = append(tse.Samples,sample)
			ts = append(ts,tse)
			req.Timeseries = ts
			fmt.Print(ts)
			ProcessReq(req)
		}
		

		
		


		
	}

}
