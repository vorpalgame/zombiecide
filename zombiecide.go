package main

import (
	"log"
	"os"
	"time"

	"github.com/vorpalgame/zombiecide/state_machines"
	"github.com/vorpalgame/core/bus"
	"github.com/vorpalgame/core/lib"
)

type zombiecide struct {
	bus bus.VorpalBus
	//textEvent  bus.TextEvent
	mouseEvent    bus.MouseEvent
	keyEvent      bus.KeyEvent
	currentZombie string
}

//Sample file for different possible use cases...

var zombies = zombiecide{}

// TODO Move to bootstrap.yaml
var fontName = "resources/font/Roboto-Regular.ttf"

func NewGame() {
	log.Println("New zombie game")

	vbus := bus.GetVorpalBus()
	//This loading/configuration needs to be consolidated into configurator
	//as we eliminate Viper in the next step.
	scene := lib.UnmarshalScene("etc/bootstrap.yaml")

	vbus.SendControlEvent(bus.NewWindowSizeEvent(scene.WindowWidth, scene.WindowHeight))
	vbus.SendControlEvent(bus.NewWindowTitleEvent(scene.WindowTitle))

	evt := bus.NewKeyRegistrationEvent(scene.RegisterKeys)
	vbus.SendKeysRegistrationEvent(evt)

	vbus.AddMouseListener(&zombies)
	vbus.AddKeyEventListener(&zombies)
	zombies.currentZombie = "h"
	zombies.bus = vbus

	zombies.mouseEvent = nil

	//MoveByIncrement to zombicide yaml

	//We're only using the first one right now...
	f, e := os.ReadFile(scene.Actors[0])
	if e != nil {
		log.Default().Println(e)
		os.Exit(1)
	}
	stateMachineZombie := state_machines.UnmarshalZombie(f)
	//log.Println(stateMachineZombie)
	//Attachable functions for testing conditions should be added so
	//they can be queried.
	//TODO we need to switch both background types to use absolute size while sprites can use percent
	//scale or perhaps both scale and width/height.
	ac := lib.ActionStageControllerData{}

	//TODO We need to revamp the configurator to eliminate Viper and to handle paths to
	//resources.
	//Need new behavior map for different environment
	ac.Load(scene.BehaviorMap)

	//TODO currently we inject this into the navigator but may
	//be better as wrapper or chain of responsiblity.
	stateMachineZombie.Navigator.ActionStageController = &ac
	textEvent := bus.NewMultilineTextEvent(fontName, 18, 0, 0).AddText("Henry follows the mouse point where legally possible.\nLeft mouse button to initiate attack!!!\nStand still too long and he dies!\n Press 'e' to exit.")
	textEvent.SetLocation(100, 100)
	//
	for {
		drawEvt := bus.NewDrawLayersEvent()
		drawEvt.AddImageLayer(*scene.Background)
		if zombies.mouseEvent != nil {
			stateMachineZombie.Execute(drawEvt, zombies.mouseEvent, zombies.keyEvent)
		}
		drawEvt.AddImageLayer(*scene.Foreground)
		vbus.SendDrawEvent(drawEvt)
		vbus.SendTextEvent(textEvent)

		zombies.keyEvent = nil
		time.Sleep(20 * time.Millisecond)

	

	}

}

func (z *zombiecide) OnKeyEvent(keyChannel <-chan bus.KeyEvent) {
	for evt := range keyChannel {
		if evt.EqualsIgnoreCase('e') {
			os.Exit(0)
		} else if evt.EqualsIgnoreCase('r') {
			//TODO Stop and close old resources if necessary...
			//zombies.zombie = NewZombie()
		} else if evt.EqualsIgnoreCase('h') {
			z.currentZombie = "h"
		} else if evt.EqualsIgnoreCase('g') {
			z.currentZombie = "g"
		} else {
			z.keyEvent = evt
		}

	}

}
func (z *zombiecide) OnMouseEvent(mouseChannel <-chan bus.MouseEvent) {
	for evt := range mouseChannel {
		z.mouseEvent = evt

	}
}
