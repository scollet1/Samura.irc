// # **************************************************************************** #
// #                                                                              #
// #                                                         :::      ::::::::    #
// #    tools.go                                           :+:      :+:    :+:    #
// #                                                     +:+ +:+         +:+      #
// #    By: scollet <marvin@42.fr>                     +#+  +:+       +#+         #
// #                                                 +#+#+#+#+#+   +#+            #
// #    Created: 2017/10/28 15:06:42 by scollet           #+#    #+#              #
// #    Updated: 2017/10/28 15:06:43 by scollet          ###   ########.fr        #
// #                                                                              #
// # **************************************************************************** #

package main

import (
  "fmt"
  "os"
)

func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
    os.Exit(1)
  }
}
