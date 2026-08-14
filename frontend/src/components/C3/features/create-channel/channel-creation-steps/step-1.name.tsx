import React, { useEffect, useState } from "react";
import { useVaultStore } from "@/store/vaultStore";
import { C3styles } from "@/components/C3/styles/styles";
import { ChannelTemplate } from "@/components/C3/domain/channel/channel.mock";
import {
    Select,
    SelectTrigger,
    SelectValue,
    SelectContent,
    SelectItem,
} from "@/components/ui/select";
import { Label } from "@radix-ui/react-label";
import { fetchChannelTemplate } from "@/components/C3/domain/channel/channel.repository";
import { CreateChannelDraft } from "../types";
import { useC3DialogStore } from "@/components/C3/infrastructure/store/c3DialogStore";


interface Step1Props {
    templates: ChannelTemplate[]
    data: CreateChannelDraft;

    onNext: (values: Partial<CreateChannelDraft>) => void;

}



export const Step1 = ({ templates, data, onNext }: Step1Props) => {
    const [channelName, setChannelName] = useState("");
    const [channelTemplate, setChannelTemplate] = useState(null)
    const [template, setTemplate] = useState<ChannelTemplate>();
    const { openC3CreateDialog, channelId } = useC3DialogStore();


    const onSelectTemplate = async (n) => {
        setChannelTemplate(n)
        data.template = n
        setChannelName(n)
        const t = await fetchChannelTemplate(n)
        setTemplate(t)
    }

    useEffect(() => {
        if (channelId) {
            onSelectTemplate(channelId as any)
        }
    }, [channelId])

    return (
        <div className="modal">
            <C3styles />
            <div className="modal-header">
                <div className="modal-title">Create Channel</div>
                <div className="step-indicator">
                    <div className="step-label">Step 1 of 4</div>
                    <div className="step-dots">
                        <div className="sdot-i active" />
                        <div className="sdot-i" />
                        <div className="sdot-i" />
                        <div className="sdot-i" />
                    </div>
                </div>
            </div>
            <div className="modal-body">

                <div className="field-label">Channel name</div>
                <div className="name-input-wrap">
                    <Select value={channelTemplate} disabled={!!channelId} onValueChange={(v) => onSelectTemplate(v as any)}>
                        <SelectTrigger>
                            <SelectValue placeholder="Select type" />
                        </SelectTrigger>
                        <SelectContent>
                            {templates.map(template => (
                                <SelectItem
                                    key={template.id}
                                    value={template.id}
                                    className="name-input"
                                >
                                    {template.title}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                    {/* <input
                        className="name-input"
                        type="text"
                        defaultValue="contract-execution"
                        readOnly
                    />
                    <span className="name-dropdown-arrow">▾</span> */}
                </div>
                <div className="templates-label">Suggested templates</div>
                <div className="template-grid">
                    <div className="tpl-card selected">
                        <div className="tpl-card-name">
                            <span className="tpl-check">✓</span>
                            Contract Execution
                        </div>
                        <div className="tpl-flow">Legal → Finance → Direction</div>
                    </div>
                    <div className="tpl-card">
                        <div className="tpl-card-name">Invoice Processing</div>
                        <div className="tpl-flow">Ops → Finance → Treasury</div>
                    </div>
                    <div className="tpl-card">
                        <div className="tpl-card-name">Budget Cycle</div>
                        <div className="tpl-flow">
                            All → Finance → Treasury → Direction
                        </div>
                    </div>
                    <div className="tpl-card">
                        <div className="tpl-card-name">Compliance Audit Close</div>
                        <div className="tpl-flow">
                            All depts → Compliance → [observer]
                        </div>
                    </div>
                    <div className="tpl-card tpl-card-single">
                        <div className="tpl-card-name">Employee Onboarding</div>
                        <div className="tpl-flow">HR → IT + Finance → Direction</div>
                    </div>
                </div>
                <div className="blank-hint">Or type any name to start blank.</div>
            </div>
            <div className="modal-footer">
                <button className="btn btn-primary" onClick={() => {
                    const selectedTpl = template || templates.find(t => t.id === channelTemplate) || templates[0];
                    const nameToUse = channelName || selectedTpl?.title || "contract-execution";
                    onNext({
                        template: selectedTpl,
                        channelName: nameToUse,
                        slots: selectedTpl?.slots ?? [],
                        properties: (selectedTpl?.defaultProperties ?? []).map(p => ({
                            key: p.key,
                            value: p.defaultValue ?? ""
                        }))
                    });
                }}>Next →</button>
            </div>
        </div>

    )
}


