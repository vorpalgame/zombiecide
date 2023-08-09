# Zombiecide

What can be more straightforward than a zombie, Henry, walking across the screen, groaning, moaning, attacking, falling
down dead only rise back up and continue on?

The first member of the zombie family is Henry, a statemachine zombie who follows the mouse pointer, groans, attacks on
command, and falls down dead if left idle too long. A newer addition is the subumption zombie, George, who is composed
of multiple separate parts that are sewn together at ever higher level to create a whole. Behavior can be overridden at
any level. Currently, for example, the head is overriden to present different ones over time. Later the leg, arms, feet
or hands may change angle or flip horizontally based inputs.

The state machine settings are read in from a yaml file, marshalled to structs, and behavior is composed by using the
names of first class functions. This provides a level of flexibiity and composability. While this is not strictly part
of the bus and engine code, it does demonstrate how the current design can accommodate those front end game designs if
desired.

The project depends on the vorpal "core" and "raylib-engine" projects

The sample code separates the zombie sprites into state machines each responsible for their own actions and transitions
to next states as well as firing events off to the engine.
![image](https://github.com/vorpalgame/vorpal/assets/3209869/95c3be51-a423-405b-8825-f5114160776d)

