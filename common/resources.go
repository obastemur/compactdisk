package resources

import "fmt"

// GetVersion returns the app version
func GetVersion() string {
    return "0.1"
}

// PrintUsage prints usage
func PrintUsage() {
  fmt.Println("compactdisk", GetVersion())
  fmt.Println("\nusage:")
  fmt.Println("compactdisk <path>")
}
