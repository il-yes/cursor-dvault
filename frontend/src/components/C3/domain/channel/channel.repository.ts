import { channel, channels, channelRows, departments, templates } from "./channel.mock"
import { ChannelView } from "./channel.types"


export const fetchChannelRows = async () => {
    return channelRows
    // return []
}   
export const fetchChannelRow = async (id: string) => {
    return channelRows.find(channel => channel.id === id)
}   
export const fetchChannel2 = async (id: string) => {
    return channel
}   
export const fetchChannel = async (id: string):Promise<ChannelView | null> => {
    return channels.find(channel => channel.id === id) || null
} 

export const fetchChannelTemplate = async (id) => {
    return templates.find( x => x.id === id)
}

export const fetchChannelTemplates = async () => {
    return templates
}

export const fetchDepartments = async () => {
    return departments
    // return []
}