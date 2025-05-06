package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/brainexploded/tortugo/config"
	"github.com/brainexploded/tortugo/inpx"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		panic(err)
	}

	inpx, err := inpx.NewInpx(cfg.LibraryPath, cfg.IndexFilename)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := inpx.Parse(ctx)
	count := 0
	var buf strings.Builder

loop:
	for {
		select {
		case inp, ok := <-ch:
			if !ok {
				break loop
			}
			if inpx.Err != nil {
				fmt.Println(inpx.Err)
				continue
			}
			buf.WriteString(fmt.Sprintf("%v %v\n", inp.Author, inp.Title))
			count++
		case <-ctx.Done():
			break loop
		default:
			continue
		}
	}
	fmt.Print(buf.String())
	fmt.Println("done, processed: ", count)
}
