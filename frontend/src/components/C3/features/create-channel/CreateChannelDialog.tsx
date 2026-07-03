// CreateChannelDialog.tsx

import { useEffect, useState } from "react";

import {
    Dialog,
    DialogContent,
} from "@/components/ui/dialog";


import { ChannelProperty, ChannelSlot, CreateChannelPayload, createEmptyChannelDraft, VaultAssignment } from "../../domain/channel/channel.types";
import { fetchChannelTemplates } from "../../domain/channel/channel.repository";
import { ChannelTemplate } from "../../domain/channel/channel.mock";
import { Step1 } from "./channel-creation-steps/step-1.name";
import { Step2 } from "./channel-creation-steps/step-2.configure";
import { Step3 } from "./channel-creation-steps/step-3.vaults";
import { C3Step4 } from "./channel-creation-steps/step-4.activate";
import { CreateChannelDraft, Props } from "./types";


export function CreateChannelDialog({ open, onClose }: Props) {
    const [templates, setTemplates] = useState<ChannelTemplate[]>([])
    const [step, setStep] = useState(1);
    const [draft, setDraft] = useState<CreateChannelDraft>(() => createEmptyChannelDraft());
    const previous = () => setStep((s) => Math.max(s - 1, 1));

    const update = (values: Partial<CreateChannelDraft>) =>
        setDraft((d) => ({
            ...d,
            ...values,
        }));


    const close = () => {

        setStep(1);

        setDraft(() => createEmptyChannelDraft());

        onClose();

    };

    const fetchTemplates = async () => {
        const templates = await fetchChannelTemplates()
        setTemplates(templates)
    }

    useEffect(() => {
        fetchTemplates()
    }, [])

    const onCreate = async(data: CreateChannelDraft) => {
        console.log({ data })
        // create channel using api
        const payload: CreateChannelPayload = {
            templateId: data.template!.id,
            title: data.channelName!,
            slots: data.slots ?? [],
            properties: data.properties ?? [],
            assignments: data.assignments ?? [],
            policy: data.policy
        }
        // update session
        close();
    }

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
                if (!o) close();
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
                    />
                )}
            </DialogContent>

        </Dialog>

    );

}