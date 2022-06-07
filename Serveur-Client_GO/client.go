package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"net"
	"os"
)

func main() {
	filepathIN := ""
	filepathOUT := ""
	for {
		//errors handling
		arguments := os.Args
		if len(arguments) == 1 {
			fmt.Println("Please provide host:port.")
			return
		}
		//connects on the provided host adress:port
		CONNECT := arguments[1]
		c, err := net.Dial("tcp", CONNECT)
		if err != nil {
			fmt.Println(err)
			return
		}

		a, _, _ := GImage(filepathIN) //getting image using filepath
		png.Encode(c, a)
		fmt.Fprintf(c, "2"+"\n")
		e, _ := png.Decode(c)

		f, _ := os.Create(filepathOUT) //create new file
		_ = png.Encode(f, e)                                             //encodes Image in new file
		f.Close()

		//infinite loop to allow client to send another request once previous one fulfilled
		/*for {
			fmt.Println("Image filepath ? or STOP to exit")
			reader := bufio.NewReader(os.Stdin) 						//requesting image filepath from user/exit message
			fmt.Print(">> ")
			text, _ := reader.ReadString('\n')

			if strings.TrimSpace(string(text)) == "STOP" { 				//exit message handling
				fmt.Println("TCP client exiting...")
				return
			}

			a,_,_ := GImage(strings.TrimSpace(string(text))) 			//getting image using filepath
			png.Encode(c, a) 											//sending image to server

			valid := false
			message, _ := bufio.NewReader(c).ReadString('\n') 	//receiving ACK response
			for !valid {
				fmt.Print("->: " + strings.TrimSpace(string(message))) 	//displaying ACK response

				fmt.Println()
				reader = bufio.NewReader(os.Stdin) 						//preparing reader
				fmt.Print(">> ")
				text, _ = reader.ReadString('\n')					//listening

				if strings.TrimSpace(text) == "1" || 					//test if response is correct
					strings.TrimSpace(text) == "2" {
					fmt.Fprintf(c, text+"\n") 							//sending response (1 or 2)
					valid = true
				} else {
					fmt.Println("Enter a valid value : 1 or 2 !")
				}
			}
			e,_:=png.Decode(c) 											//receiving image after treatment + decoding it

			fmt.Println("Output file name ?") 						//asking for output file name to user
			reader = bufio.NewReader(os.Stdin)							//preparing reader
			fmt.Print(">> ")
			text, _ = reader.ReadString('\n')						//listening
			f, _ := os.Create(strings.TrimSpace(string(text))) 			//create new file
			_ = png.Encode(f, e)										//encodes Image in new file
			f.Close()
		}*/
	}


}

/*
func GImage -> allow us to load the image using filepath
Launched from the main
Arg: string
Return: image.Image, image.Point, error
*/
func GImage(filepath string) (image.Image, image.Point, error) {
	f, err := os.Open(filepath)
	if err != nil{
		log.Fatal(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	return i, i.Bounds().Size(), err
}
