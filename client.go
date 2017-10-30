// # **************************************************************************** #
// #                                                                              #
// #                                                         :::      ::::::::    #
// #    client.go                                          :+:      :+:    :+:    #
// #                                                     +:+ +:+         +:+      #
// #    By: scollet <marvin@42.fr>                     +#+  +:+       +#+         #
// #                                                 +#+#+#+#+#+   +#+            #
// #    Created: 2017/10/28 12:46:40 by scollet           #+#    #+#              #
// #    Updated: 2017/10/28 12:46:41 by scollet          ###   ########.fr        #
// #                                                                              #
// # **************************************************************************** #

package main

import (
  "bufio"
  "net"
  "os"
  "fmt"
  "bytes"
  //"encoding/binary"
  //"io/ioutil"
  //"encoding/gob"
  //"unicode"
)

type Samurai struct {
  Nick string
  ID string
  simpleString string
  pen bytes.Buffer
  dragon bytes.Buffer
  connected bool
  nameV bool
  passV bool
  conn net.Conn
  err error
}

var Retries int = 3
const DEATH = "\u2620"

/*

  NOTE :
  to pass data correctly, the client and server must line up write/read
  operations. Double-check that there is one read for every write and
  vice versa

*/

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}

func isValid(warrior Samurai) bool {
  var finished bool
  made := false
  var input string
  var response string
  reader := bufio.NewReader(warrior.conn)
  readIn := bufio.NewReader(os.Stdin)
  for {
    /* Loop to handle new users */

    for {
      fmt.Print("New user?: [y/n] ")
      input, _ = readIn.ReadString('\n')

      if input == "y\n" || input == "yes\n" {
        _, warrior.err = warrior.conn.Write([]byte("NewUserIncoming\n"))
        for {
          fmt.Print("Username ~=> ")

          user, _ := readIn.ReadString('\n')

          //fmt.Println(user) /* NOTE */

          _, warrior.err = warrior.conn.Write([]byte(user))
          checkError(warrior.err)
          response, warrior.err = reader.ReadString('\n')
          checkError(warrior.err)

          fmt.Print(response)

          if response == "Username already registered ... \n" {
            made = true
          } else if response == "Continue ... \n" {
            warrior.ID = user
            break
          } else {
            panic("unexpected error")
            continue
          }
        }
        if made {
          break
        }
        for {
          fmt.Print("Nickname ~=> ")
          nick, _ := readIn.ReadString('\n')
          _, warrior.err = warrior.conn.Write([]byte(nick))
          checkError(warrior.err)
          warrior.Nick = nick
          break
        }
        for {
          fmt.Println(DEATH + "  »~-> \x1b[31;1mWarning: \x1b[0m" + "DO NOT USE YOUR INTRA PASSWORD <-~« " + DEATH)
          fmt.Print("Password ~=> ")
          pass, _ := readIn.ReadString('\n')
          _, warrior.err = warrior.conn.Write([]byte(pass))
          checkError(warrior.err)
          input, warrior.err = bufio.NewReader(warrior.conn).ReadString('\n')
          checkError(warrior.err)
          if input == "VALID\n" {
            fmt.Println("User created ... ")
            fmt.Println("Proceed ... ")
            made = true
          }
          if made {
            break
          } else {
            panic("something went wrong, blame the Samurai")
            continue
          }
        }
        if made {
          break
        } else {
          panic("something went wrong, blame the Samurai")
          continue
        }
      } else if input == "n\n" || input == "no\n" {
        _, warrior.err = warrior.conn.Write([]byte("ExistingUser\n"))
        break
      } else if input == "EXIT\n" {
        return false
      }  else {
        fmt.Print("Options: [y/n]")
        continue
      }
    }

    /* Loop to get user:auth */

    for {
      if (Retries == 0) {
        fmt.Println("\x1b[31;1m~*!!!!*~: \x1b[0m <!反逆者!> \x1b[31;1m~*!!!!*~: \x1b[0m")
        os.Exit(69)
      }
      fmt.Print("Username ~=> ")
      user, err := readIn.ReadString('\n')
      checkError(err)
      warrior.ID = user
      _, warrior.err = warrior.conn.Write([]byte(user))
      checkError(warrior.err)
      response, warrior.err = reader.ReadString('\n')
      checkError(warrior.err)
      if (response != "VALID\n") {
        fmt.Println("Name invalid, try again...")
        Retries -= 1
        continue
      }
      break
    }
    for {
      fmt.Print("Password ~=> ")
      pass, err := readIn.ReadString('\n')
      checkError(err)
      _, warrior.err = warrior.conn.Write([]byte(pass))
      response, warrior.err = reader.ReadString('\n')
      checkError(err)
      if (response != "VALID\n") {
        fmt.Print("Your key does not fit this lock, try again...")
        continue
      }
      break
    }
    fmt.Print("Turn the key?: [y/n] ")
    conf, _ := readIn.ReadString('\n')
    for {
      if conf == "y\n" || conf == "yes\n" {
        _, warrior.err = warrior.conn.Write([]byte("YES\n"))
        checkError(warrior.err)
        finished = true
        break
      } else if conf == "n\n" || conf == "no\n" {
        finished = false
        break
      } else {
        fmt.Println("Options: [y/n]")
        continue
      }
    }
    if finished {
      fmt.Println("returning control to main func()")
      break
    }
  }
  warrior.connected = true
  return warrior.connected
}

func main() {
  if len(os.Args) != 2 {
    fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
    os.Exit(1)
  }
  warrior := *new(Samurai)
  service := os.Args[1]
  tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
  checkError(err)
  warrior.conn, warrior.err = net.DialTCP("tcp", nil, tcpAddr)
  checkError(warrior.err)
  fmt.Println("\x1b[31;1m~*武士*~ Welcome to Samura.irc ~*チャット*~\x1b[0m" +
              "\nYou are connected to Dojo ~=> " + tcpAddr.String())
  _, warrior.err = warrior.conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
  checkError(warrior.err)
  if isValid(warrior) {
    reader := bufio.NewReader(warrior.conn)
    readIn := bufio.NewReader(os.Stdin)
    for warrior.connected {
      fmt.Print(warrior.Nick + "~*~ 武士 ~=> ")
      message, err := readIn.ReadString('\n')
      checkError(err)
      if (message == "EXIT~") {
        break
      }
      _, warrior.err = warrior.conn.Write([]byte(message))
      reply, err := reader.ReadString('\n')
      checkError(err)
      fmt.Print(reply)
    }
  }
  fmt.Println("Goodbye")
  os.Exit(0)
}
