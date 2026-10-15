package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.1"

func PrintVersion() {
	fmt.Printf("Current urlfilter version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `                                                                 
                __ ____ _  __ __             
  __  __ _____ / // __/(_)/ // /_ ___   _____
 / / / // ___// // /_ / // // __// _ \ / ___/
/ /_/ // /   / // __// // // /_ /  __// /    
\__,_//_/   /_//_/  /_//_/ \__/ \___//_/
`
	fmt.Printf("%s\n%40s\n\n", banner, "Current urlfilter version "+version)
}

