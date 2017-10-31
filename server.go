// # **************************************************************************** #
// #                                                                              #
// #                                                         :::      ::::::::    #
// #    server.go                                          :+:      :+:    :+:    #
// #                                                     +:+ +:+         +:+      #
// #    By: scollet <marvin@42.fr>                     +#+  +:+       +#+         #
// #                                                 +#+#+#+#+#+   +#+            #
// #    Created: 2017/10/28 12:46:46 by scollet           #+#    #+#              #
// #    Updated: 2017/10/28 12:46:47 by scollet          ###   ########.fr        #
// #                                                                              #
// # **************************************************************************** #

package main

import (
  "net"
  "os"
  "fmt"
  "bytes"
  "encoding/gob"
  "bufio"
  "strings"
)

const PORT = ":4242"

/*

  TODO :
  main opens the socket for connections
  start accepting connections forever
    validate users w/ the database by hashing their pass
      offer to create new user if user doesn't exist
    add validated users to runtime hashmap
    stack input from users into a runtime buffer queue
    once all inputs have been processed into the buffer
      broadcast the chat buffer to all connected clients in the hashmap
    repeat

*/

type Scribe struct {
  scribe bytes.Buffer
}

type Samurai struct {
  Nick string
  ID string
  simpleString string
  connected bool
  nameV bool
  passV bool
  conn net.Conn
  enc gob.Encoder
  dec gob.Decoder
  err error
  exists bool
  rooms map[string]*Scribe
}

type Lockbox struct {
  lock map[string]*Samurai
  Nick string
  ID string
  board bytes.Buffer
  pm bytes.Buffer
}

type Dojo map[string]*Lockbox

func saveDojo(file string, dojo map[string]*Lockbox) {
  f, err := os.Create(file)
	if err == nil {
		enc := gob.NewEncoder(f)
	  err = enc.Encode(dojo)
    checkError(err)
	}
	f.Close()
  fmt.Println("/道~ Dojo Saved ~場\\")
}

func loadDojo(file string) (dojo map[string]*Lockbox) {
  dojo = make(Dojo)
  f, err := os.Open(file)
	if err == nil {
		dec := gob.NewDecoder(f)
		err = dec.Decode(dojo)
    checkError(err)
	}
	f.Close()
  return dojo
}

func handle_connection(conn net.Conn, dojo map[string]*Lockbox) bool {
  var user string
  var pass string
  reader := bufio.NewReader(conn) /* IO reader to network */

  for {
    message, err := reader.ReadString('\n') /* Read HTTP header */
    checkError(err)

    fmt.Println("User connecting ... ")

    message, err = reader.ReadString('\n') /* Read for UserState */
    checkError(err)

    message, err = reader.ReadString('\n')

    if message == "NewUserIncoming\n" { /* Check UserState */
      user, err = reader.ReadString('\n')
      checkError(err)

      fmt.Print("Samurai \x1b[31;1mo——{=======»\x1b[0m " + user)

      if _, ok := dojo[user]; ok {
        _, err = conn.Write([]byte("Username already registered ... \n"))
        goto Validate
      } else {
        _, err = conn.Write([]byte("Continue ... \n"))

        checkError(err)
        dojo[user] = new(Lockbox)
        nick, err := reader.ReadString('\n')
        checkError(err)
        pass, err = reader.ReadString('\n')
        checkError(err)

        dojo[user].lock = make(map[string]*Samurai)
        dojo[user].lock[pass] = new(Samurai)

        if _, ok = dojo[user].lock[pass]; ok {
          dojo[user].lock[pass].Nick = nick
          _, err = conn.Write([]byte("VALID\n"))
          checkError(err)
          goto Validate
        }
      }
    } else {
      goto Validate
    }
    Validate:
      for {
        user, err = reader.ReadString('\n')
        fmt.Print("Samurai \x1b[31;1mo——{=======»\x1b[0m " + user)
        if _, ok := dojo[user]; ok {
          _, err = conn.Write([]byte("VALID\n"))
          pass, err = reader.ReadString('\n')
          if _, ok := dojo[user].lock[pass]; ok {
            _, err = conn.Write([]byte("VALID\n"))
            answer, err := reader.ReadString('\n')
            checkError(err)
            if answer == "YES\n" {
              dojo[user].lock[pass].connected = true
              message, _ := reader.ReadString('\n')
              if message == "~*~ 武士 ~=>\n" {
                goto Chat
              }
            } else {
              continue
            }
          }
        } else {
          _, err = conn.Write([]byte("No one exists under this name ... \n"))
          continue
        }
        break
      }
    Chat:
      for {
        for _, user := range dojo {
          for _, samurai := range user.lock {
            if samurai.connected {
              for _, room := range samurai.rooms {

                reader := bufio.NewReader(samurai.conn)
                scribe := new(Scribe)
                var enter string
                var pigeon string

                message, err := reader.ReadString('\n')
                fmt.Println(message)

                if (message == "NAMES~\n") {
                  for each, _ := range dojo {
                    scribe.scribe.WriteString(each)
                  }
                  samurai.conn.Write([]byte(scribe.scribe.String()))

                } else if (message == "JOIN~\n") {
                  enter, err = reader.ReadString('\n')
                  if _, ok := samurai.rooms[enter]; ok {
                    samurai.conn.Write([]byte("Room Joined\n"))
                  } else {
                    samurai.rooms[enter] = new(Scribe)
                    samurai.conn.Write([]byte("Room Created\n"))
                  }
                } else if (message == "PART~\n") {
                  leave, err := reader.ReadString('\n')
                  checkError(err)
                  if _, ok := samurai.rooms[leave]; ok {
                    delete(samurai.rooms, leave)
                  }
                } else if (message == "NICK~\n") {
                  message, err = reader.ReadString('\n')
                  samurai.Nick = message
                  samurai.conn.Write([]byte("Nickname Changed\n"))
                } else if (message == "LIST~\n") {
                  message, err = reader.ReadString('\n')
                  samurai.Nick = message
                  samurai.conn.Write([]byte(""))
                } else if (message == "PM~\n") {
                  pigeon, err = reader.ReadString('\n')
                  if _, ok := dojo[pigeon]; ok {
                    message, err = reader.ReadString('\n')
                    checkError(err)
                    dojo[pigeon].board.Write([]byte(message))
                  }
                } else {
                  room.scribe.Write([]byte(strings.TrimSuffix(samurai.Nick, "\n") + " |--- 武士 ---> " + message))
                  samurai.conn.Write([]byte(room.scribe.String()))
                }
              }
            }
          }
        }
      }
    saveDojo("dojo.dj", dojo)
    conn.Close()
    break
  }
  return false
}

func main() {
  addrs, err := net.InterfaceAddrs()
  /* NOTE : hashmap for handling runtime connections */
  dojo := loadDojo("dojo.dj") /* TODO : get the database working */
	if err != nil {
		os.Stderr.WriteString("Sh!t's F*$#ed: " + err.Error() + "\n")
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
  			os.Stdout.WriteString("Dojo running at: " + ipnet.IP.String() + "\n")
        tcpAddr, err := net.ResolveTCPAddr("tcp4", ipnet.IP.String() + PORT)
        checkError(err)
        listener, err := net.ListenTCP("tcp", tcpAddr)
        checkError(err)
        for {
          conn, err := listener.Accept() /* lock connections */
          checkError(err)
          go handle_connection(conn, dojo)
        }
      }
    }
  }
  fmt.Println("Goodbye")
  os.Exit(0)
}
