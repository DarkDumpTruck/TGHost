import axios from "axios"

const v1 = axios.create({
    baseURL: "/api/v1",
    timeout: 1000,
})

export interface Player {
    index: number
    id: number
    name: string
}

export interface Room {
    id: number
    name: string
    gameName: string
    running: boolean
    players: Player[]
}

export async function listRoom (): Promise<Room[]> {
    const response = await v1.get("/room/list")
    return response.data
}

export async function createRoom (params: {
    name: string,
    code: string,
    playerNum: number,
    botNum: number,
    judgeNum: number,
    hidden: boolean,
}) {
    const response = await v1.post("/room/create", params)
    return response.data
}

export interface PlayerStatus {
    gameName: string
    gameStatus: string
    gameRunning: boolean
    inputDone: boolean
    inputDDL: number
    inputId: string
    inputMsg: string
    inputType: string
}

export async function getPlayerStatus (params: {roomId: number, playerId: number }): Promise<PlayerStatus> {
    const response = await v1.post(`/game/status`, params)
    return response.data
}

export async function postPlayerInput (params: {roomId: number, playerId: number, inputId: string, msg: string }): Promise<{ status?: string }> {
    const response = await v1.post(`/game/input`, params)
    return response.data
}

