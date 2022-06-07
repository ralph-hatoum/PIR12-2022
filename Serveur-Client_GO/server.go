package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_"log"
	"math"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var wg2 sync.WaitGroup
var wg3 sync.WaitGroup
var wg4 sync.WaitGroup


func main() {
	//errors handling
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	//closing of connection deferred to the end
	defer l.Close()
	connum := 1
	//loop allowing us to accept potential other connections from clients after the first one and handle their requests
	for {
		fmt.Printf("#DEBUG MAIN Accepting next connection\n")
		conn, errconn := l.Accept()
		if errconn != nil {
			fmt.Printf("DEBUG MAIN Error when accepting next connection\n")
			panic(errconn)
		}
		//If we're here, we did not panic and conn is a valid handler to the new connection
		go handleConnection(conn)
		connum += 1
	}
}

/*
func handleConnection -> function handling our connection.
Launched from the main and every time there is a new connection.
Arg: net.Conn
*/
func handleConnection(c net.Conn) {
	defer c.Close()
	for{
		a, err := png.Decode(c) 										//receiving and decoding client's image
		if err!= nil{
			return
		}


		//c.Write([]byte("1 : Grey --- 2 : Gauss" + "\n")) 				//asks for the client's choice
		netData, _ := bufio.NewReader(c).ReadString('\n') 		//receives client's answer (1 or 2)
		fmt.Println("-> ", strings.TrimSpace(string(netData)))

		//execute one or the other option depending on user's choice
		//Image in Black&White : Grey

		if strings.TrimSpace(string(netData)) == "1" {
			tab := ImageToTab(a)
			tabGrey := Grey(tab)
			newImage := TabToImage(tabGrey)
			_ = png.Encode(c, newImage)
		} else if strings.TrimSpace(string(netData)) == "2" { 			//Image blurred : Gauss
			tab := ImageToTab(a)
			tabGauss := Gauss(tab)
			newImage := TabToImage(tabGauss)
			_ = png.Encode(c, newImage)
		} else {
			fmt.Println("Please give a valid entry")
		}
	}
}

/*
func ImageToTab -> function converting image to pixels array
Launched from func handleConnection before treating the image with Grey() or Gauss()
Arg: image.Image
Return: [][]color.Color
*/
func ImageToTab(im image.Image) [][]color.Color {
	numCPU := runtime.NumCPU()
	p := im.Bounds().Size()
	var tab = make([][]color.Color, p.X)
	for t := 0; t < p.X; t++ {
		tab[t] = make([]color.Color, p.Y)
	}

	chunkSize := (len(tab) + numCPU - 1) / numCPU
	for v := 0; v < len(tab); v += chunkSize {
		end := v + chunkSize
		if end > len(tab) {
			end = len(tab)
		}

		wg.Add(1)
		go func(v int, tab [][]color.Color) {
			defer wg.Done()
			for k := v; k < end; k++ {
				for j := 0; j < p.Y; j++ {
					tab[k][j] = im.At(k, j)
				}
			}
		}(v, tab)
	}

	wg.Wait()
	return tab
}

/*
func TabToImage -> function converting pixels array to image
Launched from func handleConnection after treating the image with Grey() or Gauss()
Arg: [][]color.Color
Return: image.Image
*/
func TabToImage(tab [][]color.Color) image.Image {
	rect := image.Rect(0, 0, len(tab), len(tab[0]))
	img := image.NewRGBA(rect)

	for x := 0; x < len(tab); x++ {
		wg2.Add(1)
		go func(x int, img *image.RGBA) {
			defer wg2.Done()
			for y := 0; y < len(tab[0]); y++ {
				q := tab[x]
				if q == nil {
					continue
				}
				p := tab[x][y]
				if p == nil {
					continue
				}
				original, ok := color.RGBAModel.Convert(p).(color.RGBA)
				if ok {
					img.Set(x, y, original)
				}
			}
		}(x, img)
	}

	wg2.Wait()
	return img
}

/*
func matrice7 -> returns a 3x3 matrix around the pixel given as argument
Launched from Gauss()
Arg: int, int, [][]color.Color
Return:[][]color.Color
*/
func matrice7(x int, y int, tab [][]color.Color) (m [][]color.Color) {
	var m7 = make([][]color.Color, 7)
	for i := 0; i < 7; i++ {
		var h []color.Color
		for j := 0; j < 7; j++ {
			if x-3+i >= 0 && j+y-3 >= 0 && i-3+x < len(tab) && j-3+y < len(tab[0]) {
				h = append(h, tab[x-3+i][y-3+j])
			} else {
				h = append(h, nil)
			}
		}
		m7[i] = h
	}
	return m7
}
/*
func matGauss -> returns the matrix that will be convoluted to each 3x3 matrix obtained with matrice7()
Launched from Gauss()
Arg: float64
Return:[][]float64
*/
func matGauss(sigma float64) (g [][]float64) {
	var gauss = make([][]float64, 7)
	for k := 0; k < 7; k++ {
		var y []float64
		for v := 0; v < 7; v++ {
			coef := 1 / (2 * math.Pi * math.Pow(sigma, 2)) * math.Exp(-(math.Pow(math.Abs(float64(3-k)), 2)+
				(math.Pow(math.Abs(float64(3-v)), 2)))/(2*math.Pow(sigma, 2)))
			y = append(y, coef)
		}
		gauss[k] = y
	}
	return gauss
}

/*
func Gauss -> function treating image to a blurred version
Launched from func handleConnection when treating the image with Gauss()
Arg: [][]color.Color
Return : [][]color.Color
*/
func Gauss(m [][]color.Color) (ga [][]color.Color) {
	start := time.Now()
	numCPU := runtime.NumCPU()
	m2 := make([][]color.Color, len(m))
	for i := 0; i < len(m); i++ {
		m2[i] = make([]color.Color, len(m[0]))
	}

	sigma := 10.0
	gauss := matGauss(sigma)
	chunkSize := (len(m) + numCPU - 1) / numCPU
	for i := 0; i < len(m); i += chunkSize {
		end := i + chunkSize
		if end > len(m) {
			end = len(m)
		}

		wg3.Add(1)
		/*we get for each pixel, the surrounding pixels.
		We then convolute this matrix to the matrix obtained with matGauss in order to blurr the whole picture
		*/
		go func(i int, m [][]color.Color) {
			defer wg3.Done()
			for u := i; u < end; u++ {
				for j := 0; j < len(m[0]); j++ {
					mat := matrice7(u, j, m)
					coefpdR := 0.0
					coefpdG := 0.0
					coefpdB := 0.0
					coefg := 0.0
					for k := 0; k < len(mat); k++ {
						for v := 0; v < len(mat[0]); v++ {
							if mat[k][v] != nil {
								r, g, b, _ := mat[k][v].RGBA()
								coef := gauss[k][v]
								coefg = coefg + coef
								coefpdR = coefpdR + coef*(float64(r)/257)
								coefpdG = coefpdG + coef*(float64(g)/257)
								coefpdB = coefpdB + coef*(float64(b)/257)
							} else {
								continue
							}
						}
					}
					m2[u][j] = color.RGBA{R: uint8(coefpdR / coefg), G: uint8(coefpdG / coefg), B: uint8(coefpdB / coefg), A: 255}
				}
			}
		}(i, m)
	}
	wg3.Wait()
	fmt.Println(time.Since(start))
	return m2
}

/*
func Grey -> function treating image to a B&W version
Launched from func handleConnection when treating the image with Grey()
Arg: [][]color.Color
Return : [][]color.Color
*/
func Grey(pixels [][]color.Color) (grey [][]color.Color) {
	start := time.Now()
	numCPU := runtime.NumCPU()
	xLen := len(pixels)
	yLen := len(pixels[0])
	newImage := make([][]color.Color, xLen)
	for i := 0; i < len(newImage); i++ {
		newImage[i] = make([]color.Color, yLen)
	}
	chunkSize := (xLen + numCPU - 1) / numCPU
	for i := 0; i < xLen; i += chunkSize {
		end := i + chunkSize
		if end > xLen {
			end = xLen
		}
		wg4.Add(1)
		/*
		Luminosity method :
		-> for each pixel we calculate its new value using: 0.21 R + 0.72 G + 0.07 B
		 */
		go func(i int, pixels [][]color.Color) {
			defer wg4.Done()
			for x := i; x < end; x++ {
				for y := 0; y < yLen; y++ {
					pixel := pixels[x][y]
					originalColor, ok := color.RGBAModel.Convert(pixel).(color.RGBA)
					if !ok {
						fmt.Println("type conversion went wrong")
					}
					grey := uint8(float64(originalColor.R)*0.21 + float64(originalColor.G)*0.72 + float64(originalColor.B)*0.07)
					col := color.RGBA{
						grey,
						grey,
						grey,
						originalColor.A,
					}
					newImage[x][y] = col
				}
			}
		}(i, pixels)
	}
	wg4.Wait()
	fmt.Println(time.Since(start))
	return newImage
}
