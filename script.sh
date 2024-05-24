#!/bin/bash

MAIN="./main/main.go"
KOTLIN="../../Raspred/untitled2/src/main/kotlin/"

GO_REMOVE="rm ./ssaExample/*"
GO_RUN="go run $MAIN"
MOVE="mv ./ssaExample/* $KOTLIN"
KOTLIN_REMOVE=$(ls $KOTLIN | grep -xv "Main.kt" | sed "s|.*|$KOTLIN&|" | xargs rm)


$GO_REMOVE

echo "Go generate 2 steps:"
time go run $MAIN
echo ""

$KOTLIN_REMOVE

echo "" && echo "Moving:"

time $MOVE

echo "Go fill only:"
time go run $MAIN -gen=false
echo ""

#echo "" && echo "All together:"

#$GO_REMOVE
#$KOTLIN_REMOVE

#time go run $MAIN && $MOVE

