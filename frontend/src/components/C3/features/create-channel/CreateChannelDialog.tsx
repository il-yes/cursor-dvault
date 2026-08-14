import { useEffect, useState } from "react";

import {
    Dialog,
    DialogContent,
} from "@/components/ui/dialog";

import { createEmptyChannelDraft } from "../../domain/channel/channel.types";
import { fetchChannelTemplates } from "../../domain/channel/channel.repository";
import { ChannelTemplate } from "../../domain/channel/channel.mock";
import { Step1 } from "./channel-creation-steps/step-1.name";
import { Step2 } from "./channel-creation-steps/step-2.configure";
import { Step3 } from "./channel-creation-steps/step-3.vaults";
import { C3Step4 } from "./channel-creation-steps/step-4.activate";
import { CreateChannelDraft, Props } from "./types";
import { createChannel } from "@/services/api";
import { useC3WorkspaceStore } from "@/components/C3/infrastructure/store/useC3WorkspaceStore";
import { useC3ChannelStore } from "@/components/C3/infrastructure/store/useC3ChannelStore";


export function CreateChannelDialog({ open, onClose }: Props) {
    const [templates, setTemplates] = useState<ChannelTemplate[]>([]);
    const [step, setStep] = useState(1);
    const [draft, setDraft] = useState<CreateChannelDraft>(() => createEmptyChannelDraft());
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const { activeWorkspaceId, workspaces } = useC3WorkspaceStore();
    const { addChannel } = useC3ChannelStore();

    const previous = () => setStep((s) => Math.max(s - 1, 1));

    const update = (values: Partial<CreateChannelDraft>) =>
        setDraft((d) => ({
            ...d,
            ...values,
        }));

    const close = () => {
        setStep(1);
        setError(null);
        setIsLoading(false);
        setDraft(() => createEmptyChannelDraft());
        onClose();
    };

    const fetchTemplates = async () => {
        const templates = await fetchChannelTemplates();
        setTemplates(templates);
    };

    useEffect(() => {
        fetchTemplates();
    }, []);

    const onCreate = async (data: CreateChannelDraft) => {
        setError(null);

        const workspaceId = activeWorkspaceId || (workspaces.length > 0 ? workspaces[0].id : null);
        if (!workspaceId) {
            setError("No active workspace found. Please select or create a workspace first.");
            return;
        }

        const title = data.channelName?.trim();
        if (!title) {
            setError("Channel name is required.");
            return;
        }

        setIsLoading(true);
        try {
            const createdChannel = await createChannel({
                workspace_id: workspaceId,
                title: title,
                template_id: data.template?.id || "default",
            });

            // Update C3 channel store with real server-created entity
            addChannel(createdChannel);

            close();
        } catch (err: any) {
            console.error("Failed to create channel via backend:", err);
            setError(err?.message || "An error occurred while creating the channel.");
        } finally {
            setIsLoading(false);
        }
    };

    const next = (values?: Partial<CreateChannelDraft>) => {
        if (values) {
            update(values);
        }
        setStep(s => Math.min(s + 1, 4));
    };

    return (
        <Dialog
            open={open}
            onOpenChange={(o) => {
                if (!o && !isLoading) close();
            }}
        >
            <DialogContent
                className="border-0 bg-transparent shadow-none max-w-[45rem] p-0"
            >
                {step === 1 && (
                    <Step1
                        templates={templates}
                        data={draft}
                        onNext={(values) => {
                            update(values);
                            next();
                        }}
                    />
                )}

                {step === 2 && (
                    <Step2
                        data={draft}
                        onBack={previous}
                        onNext={next}
                    />
                )}

                {step === 3 && (
                    <Step3
                        data={draft}
                        onBack={previous}
                        onNext={next}
                    />
                )}

                {step === 4 && (
                    <C3Step4
                        data={draft}
                        onBack={previous}
                        onCreate={() => onCreate(draft)}
                        isLoading={isLoading}
                        error={error}
                    />
                )}
            </DialogContent>
        </Dialog>
    );
}