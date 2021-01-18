package main

import (
	"math"
	"strconv"
)

type RunDetails struct {
	Args      []string
	Options   []string
	Page      int
	QuaverID  int
	Is7K      bool
	PageGiven bool
}

func MakeRunDetails(args []string, c chan<- RunDetails) {
	go func() {
		ctx := RunDetails{
			Args:      []string{},
			Options:   []string{},
			Page:      1,
			QuaverID:  0,
			Is7K:      false,
			PageGiven: false,
		}
		for i := 0; i < len(args); i++ {
			if i == 0 {
				continue
			}
			if args[i][0:1] == "-" {
				if len(args[i]) > 2 && args[i][0:2] == "--" {
					ctx.Options = append(ctx.Options, args[i][2:3])
				} else if len(args[i]) == 2 && args[i][1:2] == "p" {
					if len(args) >= i+1 && len(args[i+1]) > 0 && len(args[i+1]) < 4 {
						num, err := strconv.ParseInt(args[i+1], 10, 64)
						if err != nil {
							ctx.Page = 1
						} else {
							if num > 100 || num <= 0 {
								ctx.Page = 100
							} else {
								ctx.Page = int(math.Round(float64(num)))
							}
						}
						i++
						ctx.PageGiven = true
						continue
					}
				} else {
					ctx.Args = append(ctx.Args, args[i])
				}
			} else {
				ctx.Args = append(ctx.Args, args[i])
			}
		}
		c <- ctx
	}()
}
