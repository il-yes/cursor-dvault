export namespace auth {
	
	export class TokenPairs {
	    access_token: string;
	    refresh_token: string;
	    user_id: number;
	
	    static createFrom(source: any = {}) {
	        return new TokenPairs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.access_token = source["access_token"];
	        this.refresh_token = source["refresh_token"];
	        this.user_id = source["user_id"];
	    }
	}

}

export namespace blockchain {
	
	export class ChallengeRequest {
	    public_key: string;
	
	    static createFrom(source: any = {}) {
	        return new ChallengeRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.public_key = source["public_key"];
	    }
	}
	export class ChallengeResponse {
	    challenge: string;
	    expires_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ChallengeResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.challenge = source["challenge"];
	        this.expires_at = source["expires_at"];
	    }
	}
	export class SignatureVerification {
	    public_key: string;
	    challenge: string;
	    signature: string;
	
	    static createFrom(source: any = {}) {
	        return new SignatureVerification(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.public_key = source["public_key"];
	        this.challenge = source["challenge"];
	        this.signature = source["signature"];
	    }
	}

}

export namespace handlers {
	
	export class LoginRequest {
	    email: string;
	    password: string;
	    publicKey?: string;
	    signedMessage?: string;
	    signature?: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.email = source["email"];
	        this.password = source["password"];
	        this.publicKey = source["publicKey"];
	        this.signedMessage = source["signedMessage"];
	        this.signature = source["signature"];
	    }
	}
	export class LoginResponse {
	    User: models.User;
	    Vault: models.VaultPayload;
	    Tokens: auth.TokenPairs;
	
	    static createFrom(source: any = {}) {
	        return new LoginResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.User = this.convertValues(source["User"], models.User);
	        this.Vault = this.convertValues(source["Vault"], models.VaultPayload);
	        this.Tokens = this.convertValues(source["Tokens"], auth.TokenPairs);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OnBoarding {
	    user_id: string;
	    user_alias: string;
	    password: string;
	    vault_name: string;
	    role: string;
	    repo_template: string;
	    encryption_policy: string;
	    federated_providers: string[];
	
	    static createFrom(source: any = {}) {
	        return new OnBoarding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.user_alias = source["user_alias"];
	        this.password = source["password"];
	        this.vault_name = source["vault_name"];
	        this.role = source["role"];
	        this.repo_template = source["repo_template"];
	        this.encryption_policy = source["encryption_policy"];
	        this.federated_providers = source["federated_providers"];
	    }
	}
	export class OnBoardingResponse {
	    Vault: models.VaultPayload;
	    User: models.User;
	    Tokens: auth.TokenPairs;
	
	    static createFrom(source: any = {}) {
	        return new OnBoardingResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Vault = this.convertValues(source["Vault"], models.VaultPayload);
	        this.User = this.convertValues(source["User"], models.User);
	        this.Tokens = this.convertValues(source["Tokens"], auth.TokenPairs);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace models {
	
	export class CardEntry {
	    id: string;
	    entry_name: string;
	    folder_id: string;
	    type: string;
	    additionnal_note?: string;
	    custom_fields?: Record<string, string>;
	    trashed: boolean;
	    is_draft: boolean;
	    is_favorite: boolean;
	    created_at: string;
	    updated_at: string;
	    owner: string;
	    number: string;
	    expiration: string;
	    cvc: string;
	
	    static createFrom(source: any = {}) {
	        return new CardEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entry_name = source["entry_name"];
	        this.folder_id = source["folder_id"];
	        this.type = source["type"];
	        this.additionnal_note = source["additionnal_note"];
	        this.custom_fields = source["custom_fields"];
	        this.trashed = source["trashed"];
	        this.is_draft = source["is_draft"];
	        this.is_favorite = source["is_favorite"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.owner = source["owner"];
	        this.number = source["number"];
	        this.expiration = source["expiration"];
	        this.cvc = source["cvc"];
	    }
	}
	export class SSHKeyEntry {
	    id: string;
	    entry_name: string;
	    folder_id: string;
	    type: string;
	    additionnal_note?: string;
	    custom_fields?: Record<string, string>;
	    trashed: boolean;
	    is_draft: boolean;
	    is_favorite: boolean;
	    created_at: string;
	    updated_at: string;
	    private_key: string;
	    public_key: string;
	    e_fingerprint: string;
	
	    static createFrom(source: any = {}) {
	        return new SSHKeyEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entry_name = source["entry_name"];
	        this.folder_id = source["folder_id"];
	        this.type = source["type"];
	        this.additionnal_note = source["additionnal_note"];
	        this.custom_fields = source["custom_fields"];
	        this.trashed = source["trashed"];
	        this.is_draft = source["is_draft"];
	        this.is_favorite = source["is_favorite"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.private_key = source["private_key"];
	        this.public_key = source["public_key"];
	        this.e_fingerprint = source["e_fingerprint"];
	    }
	}
	export class NoteEntry {
	    id: string;
	    entry_name: string;
	    folder_id: string;
	    type: string;
	    additionnal_note?: string;
	    custom_fields?: Record<string, string>;
	    trashed: boolean;
	    is_draft: boolean;
	    is_favorite: boolean;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new NoteEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entry_name = source["entry_name"];
	        this.folder_id = source["folder_id"];
	        this.type = source["type"];
	        this.additionnal_note = source["additionnal_note"];
	        this.custom_fields = source["custom_fields"];
	        this.trashed = source["trashed"];
	        this.is_draft = source["is_draft"];
	        this.is_favorite = source["is_favorite"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class IdentityEntry {
	    id: string;
	    entry_name: string;
	    folder_id: string;
	    type: string;
	    additionnal_note?: string;
	    custom_fields?: Record<string, string>;
	    trashed: boolean;
	    is_draft: boolean;
	    is_favorite: boolean;
	    created_at: string;
	    updated_at: string;
	    genre?: string;
	    firstname?: string;
	    second_firstname?: string;
	    lastname?: string;
	    username?: string;
	    company?: string;
	    social_security_number?: string;
	    ID_number?: string;
	    driver_license?: string;
	    mail?: string;
	    telephone?: string;
	    address_one?: string;
	    address_two?: string;
	    address_three?: string;
	    city?: string;
	    state?: string;
	    postal_code?: string;
	    country?: string;
	
	    static createFrom(source: any = {}) {
	        return new IdentityEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entry_name = source["entry_name"];
	        this.folder_id = source["folder_id"];
	        this.type = source["type"];
	        this.additionnal_note = source["additionnal_note"];
	        this.custom_fields = source["custom_fields"];
	        this.trashed = source["trashed"];
	        this.is_draft = source["is_draft"];
	        this.is_favorite = source["is_favorite"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.genre = source["genre"];
	        this.firstname = source["firstname"];
	        this.second_firstname = source["second_firstname"];
	        this.lastname = source["lastname"];
	        this.username = source["username"];
	        this.company = source["company"];
	        this.social_security_number = source["social_security_number"];
	        this.ID_number = source["ID_number"];
	        this.driver_license = source["driver_license"];
	        this.mail = source["mail"];
	        this.telephone = source["telephone"];
	        this.address_one = source["address_one"];
	        this.address_two = source["address_two"];
	        this.address_three = source["address_three"];
	        this.city = source["city"];
	        this.state = source["state"];
	        this.postal_code = source["postal_code"];
	        this.country = source["country"];
	    }
	}
	export class LoginEntry {
	    id: string;
	    entry_name: string;
	    folder_id: string;
	    type: string;
	    additionnal_note?: string;
	    custom_fields?: Record<string, string>;
	    trashed: boolean;
	    is_draft: boolean;
	    is_favorite: boolean;
	    created_at: string;
	    updated_at: string;
	    user_name: string;
	    password: string;
	    web_site?: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entry_name = source["entry_name"];
	        this.folder_id = source["folder_id"];
	        this.type = source["type"];
	        this.additionnal_note = source["additionnal_note"];
	        this.custom_fields = source["custom_fields"];
	        this.trashed = source["trashed"];
	        this.is_draft = source["is_draft"];
	        this.is_favorite = source["is_favorite"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.user_name = source["user_name"];
	        this.password = source["password"];
	        this.web_site = source["web_site"];
	    }
	}
	export class Entries {
	    login: LoginEntry[];
	    card: CardEntry[];
	    identity: IdentityEntry[];
	    note: NoteEntry[];
	    sshkey: SSHKeyEntry[];
	
	    static createFrom(source: any = {}) {
	        return new Entries(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.login = this.convertValues(source["login"], LoginEntry);
	        this.card = this.convertValues(source["card"], CardEntry);
	        this.identity = this.convertValues(source["identity"], IdentityEntry);
	        this.note = this.convertValues(source["note"], NoteEntry);
	        this.sshkey = this.convertValues(source["sshkey"], SSHKeyEntry);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Folder {
	    id: number;
	    name: string;
	    created_at: string;
	    updated_at: string;
	    is_draft: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Folder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.is_draft = source["is_draft"];
	    }
	}
	
	
	
	
	export class User {
	    id: number;
	    username: string;
	    email: string;
	    password: string;
	    role: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    // Go type: time
	    last_connected_at: any;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.email = source["email"];
	        this.password = source["password"];
	        this.role = source["role"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.last_connected_at = this.convertValues(source["last_connected_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UserDTO {
	    id: number;
	    email: string;
	    role: string;
	    created_at: string;
	    updated_at: string;
	    last_connected_at: string;
	
	    static createFrom(source: any = {}) {
	        return new UserDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.email = source["email"];
	        this.role = source["role"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.last_connected_at = source["last_connected_at"];
	    }
	}
	export class VaultPayload {
	    version: string;
	    name: string;
	    folders: Folder[];
	    entries: Entries;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new VaultPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.name = source["name"];
	        this.folders = this.convertValues(source["folders"], Folder);
	        this.entries = this.convertValues(source["entries"], Entries);
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

