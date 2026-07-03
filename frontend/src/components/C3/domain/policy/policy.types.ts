// domain/policy.ts

export type AccessPolicy = {

    read:string[];

    write:string[];

    append?:string[];

    expiresAt?:string;

};