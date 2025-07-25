package display

import (
	"github.com/fatih/color"
)

func WelcomeMessage() {
	cliName := `
 ____              _____ _                    
|  _ \  _____   __/ ___|| |__   __ _ _ __ ___ 
| | | |/ _ \ \ / /\___ \| '_ \ / _ | '__/ _ \
| |_| |  __/\ V /  ___) | | | | (_| | | |  __/
|____/ \___| \_/  |____/|_| |_|\__,_|_|  \___|
`

	color.Blue(cliName)
	color.Green("DevShare - Share your dev environment over LAN")
	color.Blue("🔗 Served by DevShare - Built with 💙 by Anophel")
}
