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
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.Method)
}

const lenPath = len("/wc/")
const readLenPath = len("/rc/")

func channelHandler(w http.ResponseWriter, r *http.Request) {
	channelName := r.URL.Path[lenPath:]

	fo, err := os.OpenFile(channelName + ".txt", syscall.O_APPEND | syscall.O_CREAT , os.ModeAppend)
	                            
    
    if err != nil { panic(err) }

    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()

     // make a write buffer
    fw := bufio.NewWriter(fo)
    
    bytesToWrite := []byte("This text was written \n\n in the channel " + channelName)
    
    buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(len(bytesToWrite)))

	fmt.Fprintf(w, "%v", uint64(len(bytesToWrite)))	

    if _, err := fw.Write(buf); err != nil {
        panic(err)
    }
    
    if _, err := fw.Write(bytesToWrite); err != nil {
        panic(err)
    }

    if err = fw.Flush(); err != nil { panic(err) }
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

func readChannelHandler(w http.ResponseWriter, r *http.Request) {
	extraUrlPaths := strings.Split(r.URL.Path[readLenPath:], "/")
	if len(extraUrlPaths) == 0	{
		fmt.Fprintf(w, "No channel selected")
		return
	}

	channelName := extraUrlPaths[0]
	offset := getOffset(extraUrlPaths);

	fmt.Fprintf(w, "Channel: %s", channelName)
	fmt.Fprintf(w, "  Offset: %v", offset)

	fo, err := os.Open(channelName + ".txt")
	if err != nil { panic(err) }

	buf := make([]byte, binary.MaxVarintLen64)
	n, err := fo.ReadAt(buf, offset)
	fmt.Printf("%v", n);
	if err != nil { fmt.Fprintf(w, "offset is way off") 
		return
	}

	myInt, _ := binary.Uvarint(buf)

	fmt.Fprintf(w, "   Next offset: %v", myInt + uint64(binary.MaxVarintLen64) + uint64(offset))

	restBuf := make([]byte, myInt)
	nrest, err := fo.ReadAt(restBuf, offset + binary.MaxVarintLen64)

	fmt.Fprintf(w, "   Data: %s", string(restBuf))
	
	fmt.Printf("%v", nrest);

}


func main() {
	http.HandleFunc("/wc/", channelHandler)
	http.HandleFunc("/rc/", readChannelHandler)
    http.HandleFunc("/", handler)
    http.ListenAndServe(":9991", nil)
}