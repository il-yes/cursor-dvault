import LedgerPage from "./components/C3/ui/ledger/LedgerPage";
import ChannelPage from "./pages/C3.channel";

export type Screen =
    | { type: "ledger" }
    | { type: "channel"; channelId: string }
    | { type: "asset"; assetId: string };


export function ScreenRouter({
    screen,
    setScreen
}: {
    screen: Screen;
    setScreen: (screen: Screen) => void;
}) {


    switch (screen.type) {


        case "ledger":

            return (

                <LedgerPage

                    onOpenChannel={
                        (id) =>

                            setScreen({
                                type: "channel",
                                channelId: id
                            })

                    }

                />

            );



        case "channel":

            return (

                <ChannelPage

                    channelId={
                        screen.channelId
                    }

                />

            );



        case "asset":

            return (

                <AssetPage

                    assetId={
                        screen.assetId
                    }

                />

            );


    }

}