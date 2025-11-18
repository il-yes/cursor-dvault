import { EntrySnapshot } from "@/types/sharing";
import { VaultEntry } from "@/types/vault";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function buildEntrySnapshot(entry: VaultEntry): EntrySnapshot {
  if (!entry) {
    return {
      entry_name: "",
      type: "note"
    };
  }

  switch (entry.type) {
    case "login":
      return {
        entry_name: entry.entry_name,
        type: "login",
        user_name: entry.user_name,
        password: entry.password,
        website: entry.web_site
      };

    case "card":
      return {
        entry_name: entry.entry_name,
        type: "card",
        cardholder_name: entry.owner,
        card_number: entry.number,
        expiration: entry.expiration,
        cvv: entry.cvc
      };

    case "identity":
      return {
        entry_name: entry.entry_name,
        type: "identity",
        firstname: entry.firstname,
        lastname: entry.lastname,
        mail: entry.mail,
        telephone: entry.telephone,
        address_one: entry.address_one,
        address_two: entry.address_two,
        city: entry.city,
        state: entry.state,
        postal_code: entry.postal_code,
        country: entry.country,
        extra_fields: entry.custom_fields ?? {}
      };

    case "note":
      return {
        entry_name: entry.entry_name,
        type: "note",
        note: entry.additionnal_note
      };

    case "sshkey":
      return {
        entry_name: entry.entry_name,
        type: "sshkey",
        private_key: entry.private_key,
        public_key: entry.public_key
      };

    default:
      return {
        entry_name: entry["entry_name"],
        type: entry["type"]
      };
  }
}
