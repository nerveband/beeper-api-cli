package cmd

import (
	"fmt"
)

// Banner displays the ASCII art banner for beeper-api-cli
const banner = `
    ____  ______ ______ ____  ______ ____
   / __ )/ ____// ____// __ \/ ____// __ \
  / __  / __/  / __/  / /_/ / __/  / /_/ /
 / /_/ / /___ / /___ / ____/ /___ / _, _/
/_____/_____//_____//_/   /_____//_/ |_|

     █████╗ ██████╗ ██╗       ██████╗██╗     ██╗
    ██╔══██╗██╔══██╗██║      ██╔════╝██║     ██║
    ███████║██████╔╝██║█████╗██║     ██║     ██║
    ██╔══██║██╔═══╝ ██║╚════╝██║     ██║     ██║
    ██║  ██║██║     ██║      ╚██████╗███████╗██║
    ╚═╝  ╚═╝╚═╝     ╚═╝       ╚═════╝╚══════╝╚═╝
`

// BannerWithVersion returns the banner with version info
func BannerWithVersion() string {
	return fmt.Sprintf("%s\n    Command-line interface for Beeper Desktop API\n    Version: %s\n", banner, Version)
}

// PrintBanner prints the ASCII art banner
func PrintBanner() {
	fmt.Print(BannerWithVersion())
}
