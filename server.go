package main

import(
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"math/rand"
)

func print_binary(s []byte){
  fmt.Printf("Received b:");
  for n := 0;n < len(s);n++{
    fmt.Printf("%d,",s[n]);
  }
  fmt.Printf("\n");
}

//Game objects/players and functionality
type Player struct{
	X int 
	Y int 
	Online bool // false if offline, true if online
	Socket *websocket.Conn
}
type Cube struct{
	X int
	Y int
}
func (c *Cube) calculateNewPosition() {
	c.X = rand.Intn(480)
	c.Y = rand.Intn(380)
}
func dataGather(pos int) string{
	s := ""
	plyr := 1
	for i := 0; i < 4; i++{
		if i != pos{
			s = s+strconv.Itoa(plyr)+" " //Player
			s = s+strconv.Itoa(players[i].X)+ " "
			s = s+strconv.Itoa(players[i].Y)
			s = s+","
			plyr++
		}
	}
	s = s+"c"+" " //Cube
	s = s+strconv.Itoa(cube.X)+" "
	s = s+strconv.Itoa(cube.Y)
	return s
}
func playerInit(pos int) string{
	s := ""
	if pos == 0{
		s = s+strconv.Itoa(40)+","+strconv.Itoa(40)
	}
	if pos == 1{
		s = s+strconv.Itoa(440)+","+strconv.Itoa(40)
	}
	if pos == 2{
		s = s+strconv.Itoa(40)+","+ strconv.Itoa(340)
	}
	if pos == 3{
		s = s+strconv.Itoa(440)+","+strconv.Itoa(340)
	}
	return s
}
var players [4]Player
var cube Cube

func remoteHandler(res http.ResponseWriter, req *http.Request){
	var err error
	
	ws, err := websocket.Upgrade(res, req, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok{
		http.Error(res, "Not a websocket handshake", 400)
		return
	}else if err != nil{
		log.Println(err)
		return
	}
	
	//find empty player slot
	var slot int
	for i:=0; i<4; i++{
		if players[i].Online == false{
			slot = i
			players[i].Socket = ws //Store new connection
			players[i].Online = true //Mark as online
			break
		}
	}
	go func(){
		index := slot
		update := ""
		for{
			messageType, p, err := players[index].Socket.ReadMessage()
			if err != nil{
				return
			}
			
			//If player initiate
			if string(p) == "init"{
				//fmt.Println(string(p))
				update = playerInit(index)
				p = []byte(update)
				
				if err = players[index].Socket.WriteMessage(messageType, p); err != nil{
					return
				}
			}
			//If player move
			if string(p[0]) == "0"{
				parseResults := strings.Split(string(p), " ")
				players[index].X, _ = strconv.Atoi(parseResults[1])
				players[index].Y, _ = strconv.Atoi(parseResults[2])
				fmt.Println(string(p));
			}
			
			//If cube collision
			if string(p[0]) == "1"{
				cube.calculateNewPosition()
			}
			
			//If get server update
			if string(p[0]) == "2"{
				update = dataGather(index)
				p = []byte(update)
				
				if err = players[index].Socket.WriteMessage(messageType, p); err != nil{
					return
				}
			}
		}	
	}()
}

func main(){
	//Initial player positions
	for i:=0; i<4; i++{
		if i == 0{
			players[i].X = 40
			players[i].Y = 40
		}
		if i == 1{
			players[i].X = 40
			players[i].Y = 460
		}
		if i == 2{
			players[i].X = 560
			players[i].Y = 40
		}
		if i == 3{
			players[i].X = 560
			players[i].Y = 460
		}
	}
	
	//Initial position for cube
	cube.calculateNewPosition()
	
	http.HandleFunc("/ws", remoteHandler)
	//http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":4000", nil)
}