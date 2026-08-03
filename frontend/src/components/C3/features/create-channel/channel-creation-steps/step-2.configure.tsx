import { C3styles } from "@/components/C3/styles/styles";
import React, { useState } from "react";
import * as SwitchPrimitive from "@radix-ui/react-switch";
import { CreateChannelDraft } from "../types";

interface Step2Props {

	data: CreateChannelDraft;
	onBack: () => void;
	onNext: (values: Partial<CreateChannelDraft>) => void;

}

export const Step2 = ({ data, onBack, onNext }: Step2Props) => {
	console.log({ data })
	const [slots, setSlots] = useState(data?.slots ?? data?.template?.slots ?? [])
	const [properties, setProperties] = useState(data?.properties ?? [])


	const toggleGate = (index: number) => {
		setSlots(prev =>
			prev.map((slot, i) =>
				i === index
					? {
						...slot,
						gated: !slot.gated,
					}
					: slot
			)
		);
	};

	const updateProperty = (index: number, value: string) => {
		setProperties(prev =>
			prev.map((p, i) =>
				i === index
					? {
						...p,
						value
					}
					: p
			)
		)
	}

	return (

		<div className="modal">
			<C3styles />
			<div className="modal-header">
				<div className="modal-header-left">
					<div className="modal-title">Configure</div>
					<div className="modal-subtitle">{data.channelName}</div>
				</div>
				<div className="step-indicator">
					<div className="step-label">Step 2 of 4</div>
					<div className="step-dots">
						<div className="sdot-i done" />
						<div className="sdot-i active" />
						<div className="sdot-i" />
						<div className="sdot-i" />
					</div>
				</div>
			</div>
			<div className="modal-body">
				<div className="section-label">Property Slots</div>
				<table className="slots-table">
					<thead>
						<tr>
							<th style={{ width: "38%" }}>Slot name</th>
							<th style={{ width: "32%" }}>Commit role</th>
							<th style={{ width: "15%" }}>Gate</th>
							<th style={{ width: "15%" }}>Custom</th>
						</tr>
					</thead>
					<tbody>
						{slots?.map((slot, index) => (
							<tr key={index}>
								<td>
									<input
										className="slot-name-input"
										type="text"
										value={slot.name}
										readOnly
									/>
								</td>
								<td>
									<span className="slot-hint">{slot.role} writes</span>
								</td>
								<td>
									<div className="">
										<Switch
											className={`gate-toggle ${!slot.gated ? "off" : ""}`}
											checked={slot.gated}
											onCheckedChange={() => {
												toggleGate(index);
											}}
										/>
									</div>
								</td>
								<td>
									<div className="custom-kv">
										<input
											className="kv-input"
											type="text"
											placeholder="key"
											style={{ width: 40 }}
										/>
										<span className="kv-sep">:</span>
										<input
											className="kv-input"
											type="text"
											placeholder="val"
											style={{ width: 40 }}
										/>
									</div>
								</td>
							</tr>

						))}
					</tbody>
				</table>
				<div className="gate-note">
					Gate ● = this slot requires the previous slot to exist before it can
					be committed. Toggle to remove.
				</div>
				<div className="section-label">Channel Custom Property</div>

				{properties?.map((property, index) => (
					<div className="custom-prop-row" key={index}>
						<div className="cp-field">
							<div className="cp-label">Key</div>
							<input
								className="cp-input"
								type="text"
								defaultValue={property.key}
								readOnly
							/>
						</div>

						<div
							style={{
								display: "flex",
								alignItems: "flex-end",
								paddingBottom: 9,
								color: "#ccc",
								fontSize: 16
							}}
						>
							:
						</div>
						<div className="cp-field">
							<div className="cp-label">Value</div>
							<input
								className="cp-input"
								type="text"
								value={property.value}
								onChange={(v) => updateProperty(index, v.target.value)}
							/>
						</div>
					</div>
				))}
				<div className="cp-hint">
					e.g. counterparty: Cipla_India · fiscal_year: 2026 · project_ref:
					PRJ-042
				</div>
			</div>
			<div className="modal-footer">
				<button className="btn " onClick={() =>
					onBack()
				}>← Back</button>
				<button className="btn btn-primary" onClick={() =>
					onNext({
						slots,
						properties,
					})
				}>Next →</button>
			</div>
		</div>
	)
}


export function Switch(props: React.ComponentProps<typeof SwitchPrimitive.Root>) {
	return (
		<SwitchPrimitive.Root
			className="
                peer inline-flex h-6 w-11 shrink-0 cursor-pointer
                items-center rounded-full border-2 border-transparent
                transition-colors
                data-[state=checked]:bg-blue-600
                data-[state=unchecked]:bg-zinc-700
            "
			{...props}
		>
			<SwitchPrimitive.Thumb
				className="
                    block h-5 w-5 rounded-full bg-white shadow
                    transition-transform
                    data-[state=checked]:translate-x-5
                    data-[state=unchecked]:translate-x-0
                "
			/>
		</SwitchPrimitive.Root>
	);
}