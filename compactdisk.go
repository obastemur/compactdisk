package main

import (
  "./common"
  "./engine"
  "fmt"
  "log"
  "os"
)

func main() {
  args := os.Args[1:]
  if len(args) < 1 {
    resources.PrintUsage()
    return
  }

  files, err := fs.GetFiles(args[0])

  if err != nil {
    log.Println(err)
    return
  }

  hashmap := make(map[string]string)
  var total int64 = 0
  var duplicate int64 = 0

  for i := range files {
    var err error
    file := files[i]
    var hash string = file.SHA

    if file.IsDir() || len(hash) == 0 {
      continue
    }

    _, prs := hashmap[hash]
    if prs == false {
      hashmap[hash] = file.Path
      total += file.Size()
    } else {
      err = os.Remove(file.Path)
      if err != nil { log.Println("Failed to delete ", file.Path, ".", err); continue }
      err = os.Link(hashmap[hash], file.Path)
      if err != nil {
        log.Println("Failed to create a hardlink for ", file.Path,
                    ". However deleting the file was successful!?\nYou might want to put the original file there from ",
                    hashmap[hash], "\nerror", err);
        return;
      }
      fmt.Print("\n - ", file.Size(), " bytes disk space saved from ", file.Path, ". ")
      duplicate += file.Size()
    }
  }

  fmt.Println("\n\nTotal", len(files), "files are processed")
  fmt.Println("Unique total remaining", total / 1024, "KB")
  fmt.Println("Total disk space saved", duplicate / 1024, "KB")
}
