<script setup lang="ts">
import { BoardApi, TheChessboard, type MoveEvent, type MoveableColor } from 'vue3-chessboard'
import 'vue3-chessboard/style.css'

let board: BoardApi
let color: MoveableColor = 'black'

const handleBoardCreated = (boardApi: BoardApi) => {
  board = boardApi
  console.log(board.getFen())
}

function handleMove(move: MoveEvent) {
  console.log(move)
  console.log(board.getFen())
}
</script>

<template>
  <TheChessboard @move="handleMove" @board-created="handleBoardCreated" :player-color="color" />
</template>
