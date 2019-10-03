package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
)

type argsList struct {
	s, e, lNumber         int
	eop                   bool
	destnation, inputFile string
}

func initFlags(args *argsList) {
	flag.IntVarP(&args.s, "start", "s", -1, "The start page")
	flag.IntVarP(&args.e, "end", "e", -1, "The end page")
	flag.IntVarP(&args.lNumber, "lineOfPage", "l", 72, "The length of a page")
	flag.StringVarP(&args.destnation, "destnation", "d", "", "The destnatioon of printing")
	flag.BoolVarP(&args.eop, "endOfPage", "f", false, "Defind the end symbol of a page")
	flag.Parse()
	filename := flag.Args()
	if len(filename) == 1 {
		args.inputFile = filename[0]
	} else if len(filename) == 0 {
		args.inputFile = ""
	} else {
		fmt.Println("Too many arguments")
	}
}

func checkFlags(args *argsList) {
	if (args.s == -1) || (args.e == -1) {
		fmt.Fprintf(os.Stderr, "The start page and end page can't be empty!\n")
		os.Exit(1)
	} else if (args.s <= 0) || (args.e <= 0) {
		fmt.Fprintf(os.Stderr, "The start page and end page should be positive!\n")
		os.Exit(1)
	} else if args.s > args.e {
		fmt.Fprintf(os.Stderr, "The start page can't be bigger than the end page!\n")
		os.Exit(1)
	} else if (args.eop == true) && (args.lNumber != 72) {
		fmt.Fprintf(os.Stderr, "You can't use -f and -l together!\n")
		os.Exit(1)
	} else if args.lNumber <= 0 {
		fmt.Fprintf(os.Stderr, "The line of page can't be less than 1 !\n")
		os.Exit(1)
	} else {
		pageType := "decided by page length."
		if args.eop == true {
			pageType = "decided by the end sign /f."
		}
		dest := args.destnation
		if len(dest) == 0 {
			dest = "null"
		}
		fmt.Fprintf(os.Stderr, "startPage: %d\nendPage: %d\ninputFile: %s\npageLength: %d\npageType: %s\nprintDestation: %s\n\n",
			args.s, args.e, args.inputFile, args.lNumber, pageType, dest)
	}
}

func readFile(args *argsList) {
	var file *os.File
	var err error
	if args.inputFile != "" {
		file, err = os.Open(args.inputFile)
		defer file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file %s\n%s", args.inputFile, err)
			os.Exit(2)
		}
	} else {
		file = os.Stdin
	}

	if len(args.destnation) == 0 {
		output(os.Stdout, file, args)
	} else {
		command := exec.Command("lp", "-d"+args.destnation)
		outFile, err := command.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open Pipe!\n")
			os.Exit(2)
		}
		output(outFile, file, args)
	}
}

func output(out interface{}, in *os.File, args *argsList) {
	var pageNum int
	if args.eop {
		pageNum = 0
	} else {
		pageNum = 1
	}
	lineNum := 0
	buffer := bufio.NewReader(in)
	for {
		var pageBuf string
		var err error

		if args.eop {
			pageBuf, err = buffer.ReadString('\f')
			pageNum++
		} else {
			pageBuf, err = buffer.ReadString('\n')
			lineNum++
			if lineNum > args.lNumber {
				pageNum++
				lineNum = 1
			}
		}

		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "errors in reading file!\n")
		}

		if pageNum >= args.s && pageNum <= args.e {
			if len(args.destnation) == 0 {
				printOut, ok := out.(*os.File)
				if ok {
					fmt.Fprintf(printOut, "%s", pageBuf)
				} else {
					fmt.Fprintf(os.Stderr, "Wrong printing type!\n")
					os.Exit(3)
				}
			} else {
				printOut, ok := out.(io.WriteCloser)
				if ok {
					printOut.Write(([]byte)(pageBuf))
				} else {
					fmt.Fprintf(os.Stderr, "Wrong printing type!\n")
					os.Exit(3)
				}
			}
		}
		if err == io.EOF {
			break
		}
	}

	if pageNum < args.s {
		fmt.Fprintf(os.Stderr, "start page bigger than total pages %d, no output written\n", pageNum)
		os.Exit(4)
	} else if pageNum < args.e {
		fmt.Fprintf(os.Stderr, "end page bigger than total pages %d\n", pageNum)
		os.Exit(4)
	}
}

func main() {
	cmd := os.Args[0]
	fmt.Printf("Program Name: %s\n", cmd)
	var args argsList
	initFlags(&args)
	checkFlags(&args)
	readFile(&args)
}
