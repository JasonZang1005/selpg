package main

import (
  "fmt"
  "os"
  "flag"
  "bufio"
  "io"
  "os/exec"
)

type Args struct {
  programName string
  startPos int
  endPos int
  pageType bool
  destination string
  src string
  pageLen int
}

func main() {
    var arg Args
    arg.programName=os.Args[0]
    flag.IntVar(&arg.startPos,"s",-1,"start position")
    flag.IntVar(&arg.endPos,"e",-1,"end position")
    flag.IntVar(&arg.pageLen,"l",-1,"number of lines")
    flag.BoolVar(&arg.pageType,"f",false,"/f")
    flag.StringVar(&arg.destination,"d","","specify destination")

    flag.Parse()

    if  arg.startPos<0 || arg.endPos<0 {
      showErr("invalid start number or end number")
    }

    if arg.startPos>arg.endPos {
      showErr("start number is should be smaller than end number")
    }

    arg.src = flag.Arg(0)

    if arg.pageType==true{
      if arg.pageLen!=-1{
          showErr("only one type is allowed")
      }
    }else{
      if arg.pageLen<1{
        arg.pageLen=72
      }
    }
    if flag.NArg()==1{
      arg.src=flag.Arg(0)
    }

    if arg.src==""{
      reader:=bufio.NewReader(os.Stdin)
      if arg.pageType==true{
        page(reader,&arg)
      }else{
        line(reader,&arg)
      }
    }else{
      file,err:=os.Open(arg.src)
      reader:=bufio.NewReader(file)
      CheckErr(err)
      if arg.pageType==true{
        page(reader,&arg)
      }else{
        line(reader,&arg)
      }
    }
}



func showErr( info string) {
  fmt.Fprintf(os.Stderr,info)
  os.Exit(1)
}


func line(reader *bufio.Reader, arg *Args){
  NumOfLine :=1
  for {
    line,err:=reader.ReadString('\n')

    if NumOfLine>arg.pageLen*(arg.startPos-1)&&NumOfLine<=arg.pageLen*arg.endPos{
      if arg.destination==""{
        fmt.Println(line)
      }else{
        cmd := exec.Command("./out")
				echoInPipe, err := cmd.StdinPipe()
        CheckErr(err)
				echoInPipe.Write([]byte(line))
        fmt.Println(line)
				echoInPipe.Close()
				cmd.Stdout = os.Stdout
				cmd.Run()
      }
    }
    if err==io.EOF{
      break
    }
    NumOfLine++
  }
  if arg.startPos > NumOfLine/arg.pageLen+1 {
		fmt.Printf("Start page is greater than end page/n")
	}
	if arg.endPos > NumOfLine/arg.pageLen+1 {
		fmt.Printf("end page is greater than page num/n")
	}
}

func page(reader *bufio.Reader, arg *Args){
  NumOfPage:=1
  for{
    page,err:=reader.ReadString('\f')

    if NumOfPage>=arg.startPos&&NumOfPage<=arg.endPos{
      if arg.destination==""{
        fmt.Println(page)
      }else{
          cmd := exec.Command("./out")          // 创建命令"./out"
  				echoInPipe, err := cmd.StdinPipe()    // 打开./out的标准输入管道
  				CheckErr(err)                            // 错误检测
  				echoInPipe.Write([]byte(page + "\n")) // 向管道中写入文本
  				echoInPipe.Close()                    // 关闭管道
  				cmd.Stdout = os.Stdout                // ./out将会输出到屏幕
  				cmd.Run()
      }

      if err==io.EOF{
        break
      }
      NumOfPage++

    }
  }
  if arg.startPos > NumOfPage {
		fmt.Printf("Start page is greater than end page/n")
	}

	if arg.endPos > NumOfPage {
		fmt.Printf("end page is greater than page num/n")
	}
}

func CheckErr(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}
