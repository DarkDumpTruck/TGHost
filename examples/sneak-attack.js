//!name=偷袭
//!player=2
"use strict";

const timeLimit = 60 /* seconds */
const totalTurn = 7
const initCoin = 90
const addCoin = 10
const singleWinScore = 3
const doubleWinScore = 7
const tripleWinScore = 13

let gameState = {
	players: [],
	turn: 0,
	scores: [],
	coins: [],
	end: false,
	history: [],
	winner: -1,
}

function initAll() {
	gameState.players = getPlayers()
	for(let i = 0; i < gameState.players.length; i++) {
		gameState.scores.push(0)
		gameState.coins.push(initCoin)
	}
}

function checkGameEnd() {
	return gameState.end;
}

function startTurn() {
	gameState.turn++
	gameState.coins[0] += addCoin
	gameState.coins[1] += addCoin
}

function stepTurn() {
	if (gameState.turn == totalTurn) {
		gameState.end = true
		gameState.winner = gameState.scores[0] > gameState.scores[1] ? 0 : 1
		if (gameState.scores[0] == gameState.scores[1]) {
			gameState.winner = gameState.coins[0] > gameState.coins[1] ? 0 : 1
			if (gameState.coins[0] == gameState.coins[1]) {
				gameState.winner = -1
			}
		}
	} else {
		alertAll('第 ' + gameState.turn + ' 轮开始')
		let checkers = [];
		for(let i = 0; i < 2; i++) {
			checkers.push(numberInRangeChecker(0, gameState.coins[i], '不可输入负数！', '你没有这么多点数！'));
		}
		let inputs = getInputs(
			'请出示点数',
			'已提交点数',
			timeLimit,
			'0',
			[0, 1],
			[`slider:0:${gameState.coins[0]}`, `slider:0:${gameState.coins[1]}`],
			checkers
		)
		alertAll('第 ' + gameState.turn + ' 轮结束\n双方出示点数：' + inputs[0] + ' vs ' + inputs[1])
		gameState.history.push(inputs)
		let x = parseInt(inputs[0])
		let y = parseInt(inputs[1])
		gameState.coins[0] -= x
		gameState.coins[1] -= y
		if (x == y) {
			gameState.scores[0]++
			gameState.scores[1]++
		} else if (x >= 5 * y) {
			gameState.scores[0] += 100
			gameState.winner = 0
			gameState.end = true
		} else if (x >= 3 * y) {
			gameState.scores[0] += tripleWinScore
		} else if (x >= 2 * y) {
			gameState.scores[0] += doubleWinScore
		} else if (x > y) {
			gameState.scores[0] += singleWinScore
		} else if (x * 5 <= y) {
			gameState.scores[1] += 100
			gameState.winner = 1
			gameState.end = true
		} else if (x * 3 <= y) {
			gameState.scores[1] += tripleWinScore
		} else if (x * 2 <= y) {
			gameState.scores[1] += doubleWinScore
		} else if (x < y) {
			gameState.scores[1] += singleWinScore
		} else {
			throw 'Impossible situation'
		}
	}
}

function finalResult() {
	if (gameState.winner == -1) {
		return '游戏结束，平局'
	} else {
		return '游戏结束，玩家 ' + gameState.players[gameState.winner].Name + ' 获胜'
	}
}

function outputStatus(player) {
	let output = `第 ${gameState.turn} 轮（共 ${totalTurn} 轮）\n`
	output += `玩家：${gameState.players[0].Name}${player.Index==0?'（你）':''} vs ${gameState.players[1].Name}${1==player.Index?'（你）':''}\n`
	output += `当前分数：${gameState.scores[0]} / ${gameState.scores[1]}\n`
	output += `剩余点数：${gameState.coins[0]} / ${gameState.coins[1]}\n`
	if (gameState.history.length > 0) {
		output += `历史出价：\n`
		for (let i = 0; i < gameState.history.length; i++) {
			output += `    ${gameState.history[i][0]} , ${gameState.history[i][1]}\n`
		}
	}
	return output
}

function main() {
	initAll()
	getInputs('游戏即将开始，请输入【准备】', '已准备', 900, '准备', [0, 1], 'input', mustEqualChecker('准备'))
	while (!checkGameEnd()) {
		startTurn()
		for(let i = 0; i < gameState.players.length; i++) {
			updateStatus(i, outputStatus(gameState.players[i]))
		}
		stepTurn()
		for(let i = 0; i < gameState.players.length; i++) {
			updateStatus(i, outputStatus(gameState.players[i]))
		}
	}

	let result = finalResult()
	alertAll(result)
	appendStatusAll(result)
	return 0
}