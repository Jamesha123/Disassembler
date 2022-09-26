package main 

import(
  "fmt"
  "os"
)

var path = "addtest1_bin.txt"

func isError(err error) bool{
  if err != nil{
    fmt.Println(err.Error())
  }
  return(err != nil)
}
  
func main(){
  fmt.Println("Opening a file")
  var file, err = os.OpenFile(path, os.O_RDWR, 0644)
  if isError(err){
    return
  }
  defer file.close()
}
