import {
  INPUT_EVENT_TYPE,
  COLOR,
  Chessboard,
  BORDER_TYPE,
} from "./cm-chessboard/src/Chessboard.js";
import {
  MARKER_TYPE,
  Markers,
} from "./cm-chessboard/src/extensions/markers/Markers.js";
import {
  PROMOTION_DIALOG_RESULT_TYPE,
  PromotionDialog,
} from "./cm-chessboard/src/extensions/promotion-dialog/PromotionDialog.js";
import { Accessibility } from "./cm-chessboard/src/extensions/accessibility/Accessibility.js";
import { Chess } from "https://cdn.jsdelivr.net/npm/chess.mjs@1/src/chess.mjs/Chess.js";

const SHADOW_DOM_ENABLED = false;

class ChessBoard extends HTMLElement {
  constructor() {
    super();

    const html = `
    <div>
    <link rel="stylesheet" href="./cm-chessboard/assets/chessboard.css" />
    <link
      rel="stylesheet"
      href="./cm-chessboard/assets/extensions/markers/markers.css"
    />
    <link
      rel="stylesheet"
      href="./cm-chessboard/assets/extensions/promotion-dialog/promotion-dialog.css"
    />
      <div class="board board-large" id="board" style="width: 800px"></div>
    </div>
    `;

    let boardHtml;
    if (!SHADOW_DOM_ENABLED) {
      this.innerHTML = html;
      boardHtml = document.getElementById("board");
    } else {
      const shadow = this.attachShadow({ mode: "open" });
      shadow.innerHTML = html;
      boardHtml = shadow.getElementById("board");
    }

    const chess = new Chess();
    let color;

    function inputHandler(event) {
      console.log("inputHandler", event);
      if (event.type === INPUT_EVENT_TYPE.movingOverSquare) {
        return; // ignore this event
      }
      if (event.type !== INPUT_EVENT_TYPE.moveInputFinished) {
        event.chessboard.removeLegalMovesMarkers();
      }
      if (event.type === INPUT_EVENT_TYPE.moveInputStarted) {
        // mark legal moves
        const moves = chess.moves({
          square: event.squareFrom,
          verbose: true,
        });
        event.chessboard.addLegalMovesMarkers(moves);
        return moves.length > 0;
      } else if (event.type === INPUT_EVENT_TYPE.validateMoveInput) {
        const move = {
          from: event.squareFrom,
          to: event.squareTo,
          promotion: event.promotion,
        };
        const result = chess.move(move);
        if (result) {
          socket.send(JSON.stringify({ ...move, color: color, type: "move" }));
          event.chessboard.state.moveInputProcess.then(() => {
            // wait for the move input process has finished
            event.chessboard.setPosition(chess.fen(), true);
          });
        } else {
          // promotion?
          let possibleMoves = chess.moves({
            square: event.squareFrom,
            verbose: true,
          });
          for (const possibleMove of possibleMoves) {
            if (possibleMove.promotion && possibleMove.to === event.squareTo) {
              event.chessboard.showPromotionDialog(
                event.squareTo,
                color === "white" ? COLOR.white : COLOR.black,
                (result) => {
                  console.log("promotion result", result);
                  if (
                    result.type === PROMOTION_DIALOG_RESULT_TYPE.pieceSelected
                  ) {
                    const move = {
                      from: event.squareFrom,
                      to: event.squareTo,
                      promotion: result.piece.charAt(1),
                    };
                    chess.move(move);
                    socket.send(
                      JSON.stringify({ ...move, color: color, type: "move" })
                    );
                    event.chessboard.setPosition(chess.fen(), true);
                  } else {
                    // promotion canceled
                    event.chessboard.enableMoveInput(inputHandler, COLOR.white);
                    event.chessboard.setPosition(chess.fen(), true);
                  }
                }
              );
              return true;
            }
          }
        }
        return result;
      } else if (event.type === INPUT_EVENT_TYPE.moveInputFinished) {
        if (event.legalMove) {
          event.chessboard.disableMoveInput();
        }
      }
    }

    const board = new Chessboard(boardHtml, {
      position: chess.fen(),
      assetsUrl: "./cm-chessboard/assets/",
      style: {
        borderType: BORDER_TYPE.none,
        pieces: { file: "pieces/staunty.svg" },
        animationDuration: 300,
      },
      orientation: COLOR.white,
      extensions: [
        { class: Markers, props: { autoMarkers: MARKER_TYPE.square } },
        { class: PromotionDialog },
        { class: Accessibility, props: { visuallyHidden: true } },
      ],
    });

    const socket = new WebSocket("ws://localhost:5555/ws");
    socket.addEventListener("message", (event) => {
      console.log("Message from server ", event.data);
      let message = JSON.parse(event.data);
      if (message.type === "start") {
        color = message.color;
        if (color === "white") {
          board.enableMoveInput(inputHandler, COLOR.white);
        } else {
          board.setOrientation(COLOR.black);
        }
      } else if (message.type === "move") {
        let move = {
          from: message.from,
          to: message.to,
          promotion: message.promotion,
        };
        let result = chess.move(move);
        console.log(result);
        board.setPosition(chess.fen(), true);
        board.enableMoveInput(
          inputHandler,
          color === "white" ? COLOR.white : COLOR.black
        );
      }
    });
  }
}

customElements.define("chess-board", ChessBoard);
