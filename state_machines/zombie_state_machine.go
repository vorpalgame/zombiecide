package state_machines

import (
	"fmt"
	"github.com/vorpalgame/core/bus"
	"github.com/vorpalgame/core/lib"
)

///////////////////////////////////////////////////////////////////////////////////////
/////      Zombie Data and Functions

type Zombie interface {
	Execute(drawEvent bus.DrawLayersEvent, mouseEvent bus.MouseEvent, keyEvent bus.KeyEvent)
	getNavigator() lib.Navigator
}

func (z *ZombieData) Execute(drawEvent bus.DrawLayersEvent, mouseEvent bus.MouseEvent, keyEvent bus.KeyEvent) {
	state := z.StateMap[z.CurrentStateName]

	tx := NewBehaviorTransaction(z.CurrentStateName, mouseEvent, z.getNavigator(), state.getFrameTracker())
	state.getBehaviors().ExecuteBehaviors(tx)

	if tx.IsTransition() {
		bus.GetVorpalBus().SendAudioEvent(bus.NewStopAudioEvent(state.AudioState))
		state.FrameTracker.ResetFrameCount()
		z.CurrentStateName = tx.GetNextStateName()
		state = z.StateMap[z.CurrentStateName]
		bus.GetVorpalBus().SendAudioEvent(bus.NewPlayAudioEvent(state.AudioState))
	}
	state.render(drawEvent, mouseEvent, z.getNavigator())
}

func (z *ZombieData) getNavigator() lib.Navigator {
	return z.Navigator
}

///////////////////////////////////////////////////////////////////////////////////////
/////          State Data and Functions

type ZombieState interface {
	getBehaviors() BehaviorsExecutor
	getFrameTracker() lib.FrameTracker
	render(drawEvent bus.DrawLayersEvent, mouseEvent bus.MouseEvent, navigator lib.Navigator)
}

func (z *ZombieStateData) getFrameTracker() lib.FrameTracker {
	//log.Default().Println(z)
	return z.FrameTracker
}

func (z *ZombieStateData) getBehaviors() BehaviorsExecutor {
	if z.behaviorsExecutor == nil {
		z.behaviorsExecutor = NewBehaviorExecutor(z.BehaviorNames)
	}
	return z.behaviorsExecutor
}

func (z *ZombieStateData) render(drawEvent bus.DrawLayersEvent, mouseEvent bus.MouseEvent, navigator lib.Navigator) {
	img := getZombieImage(z.Spec, z.Name, z.FrameTracker.GetCurrentFrame())
	x, y := navigator.GetCurrentPoint()
	metadata := lib.ImageMetadata{img, x, y, z.Width, z.Height, flipHorizontal(mouseEvent, navigator)}
	layer := lib.ImageLayerData{}
	layer.LayerMetadata = append(layer.LayerMetadata, &metadata)
	drawEvent.AddImageLayer(layer)
}

// Helper functions...
func flipHorizontal(mouseEvent bus.MouseEvent, locator lib.Navigator) bool {
	mx := mouseEvent.GetCursorPoint().GetX()
	lx, _ := locator.GetCurrentPoint()
	return mx < lx
}
func getZombieImage(spec, name string, frameNumber int32) string {
	return fmt.Sprintf(spec, name, frameNumber)
}
