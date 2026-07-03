import { ThreadAssetViewInterface } from "../../domain/thread/asset.types";


export const ThreadAssetHeader = ({ asset }: { asset: ThreadAssetViewInterface }) => {
    return (
        <div className="thread-asset-header">
            <div className="thread-asset-header-title">
                <h2>{asset.title}</h2>
                <p>{asset.subtitle}</p>
            </div>
            <div className="thread-asset-header-actions">
                <button>Timeline</button>
                <button>Receipts</button>
                <button>Policies</button>
                <button>Anchors</button>
            </div>
        </div>
    );
}
