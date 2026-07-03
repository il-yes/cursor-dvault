import { threadAssets } from "./asset.mock"


export const fetchThreadAsset = async (assetId: string) => {
    return threadAssets.find(x => x.id === assetId) 
}
