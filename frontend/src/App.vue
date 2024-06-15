<script setup lang="ts">
import { ref } from 'vue'
import { BoardApi, TheChessboard, type MoveEvent } from 'vue3-chessboard'
import 'vue3-chessboard/style.css'

let board: BoardApi
const color = ref()
const whiteTime = ref()
const blackTime = ref()

const socket = new WebSocket('ws://localhost:5555/ws')
socket.addEventListener('message', (event) => {
  const message = JSON.parse(event.data)
  if (message.type === 'start') {
    color.value = message.color
    whiteTime.value = message.whiteTime
    blackTime.value = message.blackTime
  } else if (message.type === 'move') {
    const { from, to, promotion } = message
    whiteTime.value = message.whiteTime
    blackTime.value = message.blackTime
    if(color.value!=message.color){
      board.move({ from, to, promotion })
    }
  }
})

function formatSeconds(seconds: number){
  const zeroPad = (num: number, places: number) => String(num).padStart(places, '0')
  const minutes = Math.floor(seconds/60)
  seconds = seconds%60
  return `${zeroPad(minutes,2)}:${zeroPad(seconds,2)}`
}

function handleBoardCreated(boardApi: BoardApi) {
  board = boardApi
}

function handleMove(move: MoveEvent) {
  if (!color.value.startsWith(move.color)) {
    return
  }
  const { from, to, promotion } = move
  const message = JSON.stringify({ from, to, promotion, color: color.value, type: 'move' })
  socket.send(message)
}

let interval = setInterval(() => {
    if(whiteTime.value === 0 || blackTime.value===0) clearInterval(interval)
    if(board.getTurnColor()==="white") whiteTime.value=whiteTime.value-1;
    if(board.getTurnColor()=="black") blackTime.value--;
    }, 1000)
</script>

<template>
  <div v-if="color">
    <h1 v-if="color === 'black'">{{formatSeconds(whiteTime)}}</h1>
    <h1 v-else>{{formatSeconds(blackTime)}}</h1>
    <TheChessboard
      @move="handleMove"
      @board-created="handleBoardCreated"
      :player-color="color"
      :board-config="{ orientation: color }"
    />
    <h1 v-if="color === 'white'">{{formatSeconds(whiteTime)}}</h1>
    <h1 v-else>{{formatSeconds(blackTime)}}</h1>
  </div>
  <h1 v-else>Waiting for player 2</h1>
</template>
