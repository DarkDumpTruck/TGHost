//!name=动物园捐赠
//!player=10
"use strict";

const timeLimitPrepare = 1999 /* seconds */
const timeLimitReplaceCard = 60 /* seconds */
const timeLimitChooseCard = 180 /* seconds */
const timeLimitBet = 60 /* seconds */
const timeLimitCall = 60 /* seconds */
const totalTurn = 10
const totalPlayer = 10
const allPlayerIds = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
const initCoin = 100

let gameState = {
	players: [],
	turn: 0,
	cards: [],
	coins: [],
	goods: [],
	stage: 0,
	showhand: false,
	winner: -1,
	endStage: false,
	showCards: [],
	backCards: [],
	currentBets: [],
	fold: [],
	history: [],
}

let allCards = []
for(let i of ['鲜花', '蜜桃', '钻石', '黑桃']) {
	for(let j of ['2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A']) {
		allCards.push(i+j)
	}
}

function initAll() {
	gameState.players = getPlayers()
	for(let i = 0; i < totalPlayer; i++) {
		gameState.coins.push(initCoin)
		gameState.goods.push(0)
		gameState.history.push("")
	}
}

function checkGameEnd() {
	for(let i = 0; i < totalPlayer; i++) {
		if(gameState.goods[i] >= 20) {
			return true
		}
	}
	return gameState.turn >= totalTurn
}

function finalResult() {
	return '游戏结束，请等待裁判公布游戏结果。'
}

function shuffle(array) {
	for(let i = 0; i < 5; i++) {
		let currentIndex = array.length,  randomIndex
		while (currentIndex > 0) {
		randomIndex = Math.floor(Math.random() * currentIndex)
		currentIndex--
		[array[currentIndex], array[randomIndex]] = [
			array[randomIndex], array[currentIndex]]
		}
	}
	return array
}

function getShuffledCard() {
	let cards = allCards.slice()
	shuffle(cards)
	return cards
}

let specialPlayerIds = shuffle(allPlayerIds.slice())

function outputStatus(index) {
	let output = ''
	if(gameState.endStage) {
		output = '第' + gameState.turn + '轮结束。\n'
		output += '场上展示的牌是：' + gameState.showCards.join(' ') + '\n'
		if(gameState.showhand) {
			output += '未展示的五张牌是：' + gameState.backCards.join(' ') + '\n\n'
			output += '玩家开牌情况：\n'
			for(let i = 0; i < totalPlayer; i++) {
				if(!gameState.fold[i]) {
					output += `${gameState.players[i].Name}：${gameState.cards[i].join(' ')}\n`
				}
			}
		} else {
			output += `因其他玩家弃牌，${gameState.players[gameState.winner].Name}获得本轮胜利。`
		}
		output += '\n各玩家剩余点数情况：\n'
		for(let i = 0; i < totalPlayer; i++) {
			output += `${gameState.players[i].Name}：${gameState.coins[i]}水晶，${gameState.goods[i]}好人卡；\n`
		}
		gameState.history[index] += output + '\n'
		return gameState.history[index]
	}
	
	if(gameState.stage == 0) {
		output = '当前第' + gameState.turn + '轮的开始阶段，你目前有' + gameState.coins[index] + '水晶，' + gameState.goods[index] + '好人卡。\n'
		output += '你目前的手牌是：' + gameState.cards[index].join(' ') + '\n\n'
	}
	else if(gameState.stage % 2 == 1) {
		output = `当前第 ${gameState.turn} 轮的第 ${(gameState.stage + 1) / 2} 阶段，你目前有` + gameState.coins[index] + '水晶，' + gameState.goods[index] + '好人卡。\n'
		output += '你目前的手牌是：' + gameState.cards[index].join(' ') + '\n'
		output += '场上展示的牌是：' + gameState.showCards.join(' ') + '\n\n'
		if(gameState.stage > 1) {
			output += '当前募捐情况：\n'
			for(let i = 0; i < totalPlayer; i++) {
				if(gameState.currentBets[i] > 0) {
					output += gameState.players[i].Name + '募捐了' + gameState.currentBets[i] + '个水晶' + (gameState.coins[i] == 0?'（全下）\n':gameState.fold[i]?'（已弃牌）\n':'。\n')
				}
			}
		}
	}
	else if(gameState.stage % 2 == 0) {
		output = `当前第 ${gameState.turn} 轮的第 ${gameState.stage / 2} 阶段，你目前有` + gameState.coins[index] + '水晶，' + gameState.goods[index] + '好人卡。\n'
		output += '你目前的手牌是：' + gameState.cards[index].join(' ') + '\n'
		output += '场上展示的牌是：' + gameState.showCards.join(' ') + '\n\n'
		output += '当前募捐情况：\n'
		for(let i = 0; i < totalPlayer; i++) {
			if(gameState.currentBets[i] > 0) {
				output += gameState.players[i].Name + '募捐了' + gameState.currentBets[i] + '个水晶' + (gameState.coins[i] == 0?'（全下）\n':gameState.fold[i]?'（已弃牌）\n':'。\n')
			}
		}
	}

	return gameState.history[index] + output
}

function updateAll() {
	for(let i = 0; i < totalPlayer; i++) {
		updateStatus(i, outputStatus(i))
	}
}

function main() {
	initAll()
	getInputs(
        '游戏即将开始，请输入【准备】',
        '已准备',
        timeLimitPrepare,
        '准备',
        allPlayerIds,
        'input',
        mustEqualChecker('准备')
    )
	while (!checkGameEnd()) {
		gameState.turn++
		let cards = getShuffledCard()
		let index = 0
		gameState.stage = 0
		gameState.showhand = false
		gameState.endStage = false
		gameState.cards = []
		gameState.fold = []
		for(let i = 0; i < totalPlayer; i++) {
			gameState.fold.push(gameState.coins[i] == 0)
		}
		
		let specialPlayerId = specialPlayerIds[gameState.turn - 1]
		for(let i = 0; i < totalPlayer; i++) {
			if(gameState.fold[i]) {
				gameState.cards.push([])
				continue
			}
			let player_cards = []
			for(let j = 0; j < 3; j++) {
				player_cards.push(cards[index++])
			}
			if(i == specialPlayerId) {
				player_cards[2] = '【特殊牌】'
			}
			gameState.cards.push(player_cards)
		}
		updateAll()
		let input = getInputs(
			'请将【特殊牌】替换为任意一张牌，超时默认替换为'+gameState.cards[specialPlayerId][0],
			'替换成功',
			timeLimitChooseCard,
			gameState.cards[specialPlayerId][0],
			[specialPlayerId],
			'input',
			itemInListChecker(allCards)
		)
		gameState.cards[specialPlayerId][2] = input[0]
		updateStatus(specialPlayerId, outputStatus(specialPlayerId))

		let survivePlayerIds = []
		for(let i = 0; i < totalPlayer; i++) {
			if(!gameState.fold[i]) {
				survivePlayerIds.push(i)
			}
		}
		let inputs = getInputs(
			'请选择一张手牌用于【物资储备】，超时默认第一张',
			'选择成功',
			timeLimitChooseCard,
			'1',
			survivePlayerIds,
			'input', // TODO: use button
			numberInRangeChecker(1, 3)
		)
		let discardCards = []
		for(let i = 0; i < totalPlayer; i++) {
			if(gameState.fold[i]) {
				continue
			}
			let cardIndex = parseInt(inputs[i]) - 1
			discardCards.push(gameState.cards[i][cardIndex])
			gameState.cards[i].splice(cardIndex, 1)
		}
		shuffle(discardCards)
		gameState.backCards = discardCards.slice(5)
		gameState.showCards = []
		for(let i = 0; i < 2; i++) {
			gameState.showCards.push(discardCards[i])
		}
		gameState.currentBets = []
		for(let i = 0; i < totalPlayer; i++) {
			gameState.currentBets.push(0)
		}
		let currentShowCardIndex = 2
		while(gameState.stage < 6) {
			gameState.stage++
			gameState.showCards.push(discardCards[currentShowCardIndex++])
			updateAll()
			let minBets = []
			let inputSliders = []
			let checkers = []
			let minBetNum = gameState.stage == 1 ? gameState.turn : 0
			for(let i of survivePlayerIds) {
				if(gameState.coins[i] >= minBetNum) {
					minBets.push(minBetNum)
					inputSliders.push('slider:' + minBetNum + ':' + gameState.coins[i])
					checkers.push(numberInRangeChecker(minBetNum, gameState.coins[i]))
				} else {
					minBets.push(gameState.coins[i])
					inputSliders.push('slider:' + gameState.coins[i] + ':' + gameState.coins[i])
					checkers.push(numberInRangeChecker(gameState.coins[i], gameState.coins[i]))
				}
			}
			inputs = getInputs(
				'请选择募捐水晶数量，本轮最低募捐数量为'+minBetNum,
				'募捐成功',
				timeLimitBet,
				minBets,
				survivePlayerIds,
				inputSliders,
				checkers
			)
			gameState.stage++
			let maxBet = 0
			for(let i = 0; i < survivePlayerIds.length; i++) {
				let index = survivePlayerIds[i]
				let bet = parseInt(inputs[i])
				gameState.currentBets[index] += bet
				gameState.coins[index] -= bet
				if(gameState.currentBets[index] > maxBet) {
					maxBet = gameState.currentBets[index]
				}
			}
			updateAll()
			let needCallIds = []
			for(let i of survivePlayerIds) {
				if(gameState.currentBets[i] < maxBet && gameState.coins[i] > 0) {
					needCallIds.push(i)
				}
			}
			inputs = getInputs(
				'请选择是否跟投，超时默认放弃',
				'选择成功',
				timeLimitCall,
				'放弃',
				needCallIds,
				'input',
				mustInChecker(['跟投', '放弃'])
			)
			for(let i = 0; i < needCallIds.length; i++) {
				let index = needCallIds[i]
				if(inputs[i] == '跟投') {
					let callNum = Math.min(gameState.coins[index], maxBet - gameState.currentBets[index])
					gameState.coins[index] -= callNum
					gameState.currentBets[index] += callNum
				} else {
					gameState.fold[index] = true
				}
			}
			let canBetCount = 0
			for(let i = 0; i < totalPlayer; i++) {
				if(gameState.fold[i] || gameState.coins[i] == 0) {
					continue
				}
				canBetCount++
			}
			if(canBetCount <= 1) {
				break
			}
		}

		let unfoldedCount = 0
		gameState.winner = -1
		for(let i = 0; i < totalPlayer; i++) {
			if(gameState.fold[i]) {
				continue
			}
			unfoldedCount++
			gameState.winner = i
		}
		if(unfoldedCount == 1) {
			let totalBet = 0
			for(let i = 0; i < totalPlayer; i++) {
				totalBet += gameState.currentBets[i]
			}
			gameState.coins[gameState.winner] += totalBet
		} else {
			gameState.showhand = true
		}
		gameState.endStage = true
		updateAll()
	}

	let result = finalResult()
	alertAll(result)
	appendStatusAll(result)
	return 0
}

