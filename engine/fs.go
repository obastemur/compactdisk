package fs

import (
  "fmt"
  "errors"
  "os"
  "io"
  "runtime"
  "math"
  "path/filepath"
  "encoding/hex"
  "crypto/sha1"
  "../common"
)

const sizeMB4 = 1024 * 1024 * 4
var threadPool resources.ThreadPool

// FileEntry structure that is extended on top of os.FileInfo
type FileEntry struct {
  os.FileInfo
  Path string
  SHA  string
  Error error
}

// GetHash compiles the sha1 hash of the file
func getHash(entry FileEntry, entries chan<- FileEntry) {
  if entry.IsDir() {
    errorString := fmt.Sprintf("Path is not a file : %s", entry.Path)
    entry.Error = errors.New(errorString)
  } else {
    file, err := os.Open(entry.Path)
    if err != nil {
      entry.Error = err
    } else {
      defer file.Close()
      hash := sha1.New()

      tmpBuffer := make([]byte, sizeMB4)
      var currentPosition int64 = 0
      var threshholdCounter int64 = 0
      isEmpty := true
      entry.SHA = ""

      for {
        readCount, err := file.Read(tmpBuffer)
        if err != nil || readCount == 0 {
          if err != io.EOF {
            entry.Error = err
          }
          break
        }

        hash.Write(tmpBuffer[:readCount])
        currentPosition += int64(readCount)
        threshholdCounter += int64(readCount)
        isEmpty = false

        _, err = file.Seek(currentPosition, 0)
        if threshholdCounter >= 256 * sizeMB4 {
          sha := hex.EncodeToString(hash.Sum(nil))
          hash = sha1.New()
          entry.SHA += sha
          isEmpty = true
          threshholdCounter = 0
        }
      }

      if entry.Error == nil && !isEmpty {
        sha := hex.EncodeToString(hash.Sum(nil))
        entry.SHA += sha
      }
    }
  }
  entries <- entry
  threadPool.Done()
}

// GetFiles returns the files in the path recursively
func GetFiles(path string) ([]FileEntry, error) {
  var files []FileEntry
  var fileCount int64 = 0

  // math min is safe as the range of comparison within the float64
  maxThreadCount := int32( math.Min(4, float64(runtime.NumCPU())) )
  threadPool = resources.NewThreadPool(maxThreadCount)

  // select all the `files` under the path recursively and resolve their fullpath
  err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    fullPath := path
    if !filepath.IsAbs(fullPath) {
      absPath, err := filepath.Abs(fullPath)
      if err != nil {
        return err
      }

      fullPath = absPath
    }

    file := FileEntry{info, fullPath, "", nil}
    files = append(files, file)

    if !file.IsDir() { fileCount++; }
    if fileCount % 111 == 0 { fmt.Print(".") }

    return nil
  })

  if err != nil {
    return nil, err
  }

  entries := make(chan FileEntry, fileCount)

  // process their SHA1
  for i := range files {
    file := files[i]

    if !file.IsDir() {
      threadPool.Add(1)
      go getHash(file, entries)
      threadPool.WaitPool()
      if i % 11 == 0 {
        fmt.Print(".")
      }
    }
  }

  threadPool.Wait()
  close(entries)

  newfiles := make([]FileEntry, 0)

  for entry := range entries {
    newfiles = append(newfiles, entry)
  }

  files = nil

  return newfiles, nil
}
