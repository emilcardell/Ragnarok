package main

import (
	"fmt"
	"net/http"
	"bufio"
	"syscall"
    "os"
    "encoding/binary"
    "strings"
    "strconv"
    "io"
    "io/ioutil"
)

const lenPath = len("/WriteToChannel/")
const readLenPath = len("/ReadFromChannel/")
const channelPath = "data/channels/"


func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Ragnarok is online. \n\nChannels: \n")

    files, _ := ioutil.ReadDir(channelPath) 
    for _, value := range files {
    	channelName := strings.Split(value.Name(), ".")[0]
    	fmt.Fprintf(w, "%s \n", channelName)
	}
}



func writeToChannelHandler(w http.ResponseWriter, r *http.Request) {
	channelName := r.URL.Path[lenPath:]

	fo, err := os.OpenFile(channelPath + channelName + ".channel", syscall.O_APPEND | syscall.O_CREAT , os.ModeAppend)
	fi, _ := fo.Stat()

	startOffset := fi.Size()
	
    if err != nil { panic(err) }

    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()

    fw := bufio.NewWriter(fo)
    


    bytesToWrite := make([]byte, r.ContentLength)
    fmt.Println("ContentLength: %v", r.ContentLength)
	numbOfBodyBytes, err := io.ReadFull(r.Body, bytesToWrite)
	fmt.Println("numbOfBodyBytes: %v", numbOfBodyBytes)
    
    buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(len(bytesToWrite)))

    if _, err := fw.Write(buf); err != nil {
        panic(err)
    }
    
    if _, err := fw.Write(bytesToWrite); err != nil {
        panic(err)
    }

    if err = fw.Flush(); err != nil { panic(err) }

    fmt.Fprintf(w, "%v", startOffset)	
}

func getOffset(urlParts []string) int64 {
	if( len(urlParts) == 2) {
		myInt, err := strconv.ParseInt(urlParts[1], 10, 64)
		
		if err != nil { 
			return int64(0);
		}

		return myInt
	}

	return 0
}

func readFromChannelHandler(w http.ResponseWriter, r *http.Request) {
	extraUrlPaths := strings.Split(r.URL.Path[readLenPath:], "/")
	if len(extraUrlPaths) == 0	{
		fmt.Fprintf(w, "No channel selected")
		return
	}

	channelName := extraUrlPaths[0]
	offset := getOffset(extraUrlPaths);

	fo, err := os.Open(channelPath + channelName + ".channel")
	if err != nil { panic(err) }

	buf := make([]byte, binary.MaxVarintLen64)
	n, err := fo.ReadAt(buf, offset)
	fmt.Println("n: %v", n)
	if err != nil { 
		fmt.Fprintf(w, "offset is way off") 
		return
	}

	myInt, _ := binary.Uvarint(buf)

	restBuf := make([]byte, myInt)
	nrest, err := fo.ReadAt(restBuf, offset + binary.MaxVarintLen64)
	fmt.Println("nrest: %v", nrest)

	fmt.Fprintf(w, "%v", restBuf)
}


func main() {
	http.HandleFunc("/WriteToChannel/", writeToChannelHandler)
	http.HandleFunc("/ReadFromChannel/", readFromChannelHandler)
    http.HandleFunc("/", handler)

    err := os.MkdirAll(channelPath, os.ModeDir)
    if err != nil { 
    	fmt.Println("Data dir could not be created") 		
	}
    fmt.Println("Ragnarok channels are running on 9991")
    http.ListenAndServe(":9991", nil)

}