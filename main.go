package main

// socksie is a SOCKS4/5 compatible proxy that forwards connections via
// SSH to a remote host

import (
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/rob121/pman/icon"
	"github.com/rob121/vhelp"
	"github.com/spf13/viper"
	"log"
	"os/exec"
)
//windows ubild env GO111MODULE=on go build -ldflags "-H=windowsgui"

 var config string
 var v *viper.Viper
 var devices []string

type Item struct{
	Device string
	Menu *systray.MenuItem
}

var items map[string]Item

func main() {

	items = make(map[string]Item)
    
    flag.StringVar(&config,"config","config","Config File Name")
    
    vhelp.Load(config)

    var verr error
    
    v,verr = vhelp.Get(config)

    if(verr!=nil){}
    
    devices = v.GetStringSlice("devices")
    
    systray.Run(onReady, onExit)
    
}

func cmd(device string, state string){

	cmd := exec.Command("networksetup","-setsocksfirewallproxystate", device, state)

	fmt.Printf("Got %s:%s\n",device,state)

	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	//fmt.Printf("combined out:\n%s\n", string(out))


}

func onReady() {


	systray.SetIcon(icon.Data)
	systray.SetTitle(fmt.Sprintf("Proxy Manager"))
	//systray.SetTooltip(fmt.Sprintf("Connected to %s connect on %s",addrs,addr))

	for _,d := range devices {

		mChecked := systray.AddMenuItemCheckbox(fmt.Sprintf("Proxy On: %s",d), "Click to activate/deactivate", false)


		items[d] = Item{d,mChecked}

		cmd(d,"off")

	}



	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	//mQuit.SetIcon(icon.Data)

	for id,_ := range items {

		go func(id string) {

			mChecked := items[id].Menu

			for {
				select {
				case <-mChecked.ClickedCh:
					if mChecked.Checked() {
						mChecked.Uncheck()

                         cmd(id,"off")
					} else {
						cmd(id,"on")
						mChecked.Check()
					}
				}
			}

		}(id)

	}
	
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
	
}

func onExit() {
	// clean up here
}
