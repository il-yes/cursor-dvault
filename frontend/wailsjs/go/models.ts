export namespace app_config {
	
	export class IPFSConfig {
	    api_endpoint: string;
	    gateway_url: string;
	
	    static createFrom(source: any = {}) {
	        return new IPFSConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.api_endpoint = source["api_endpoint"];
	        this.gateway_url = source["gateway_url"];
	    }
	}
	export class StellarConfig {
	    network: string;
	    horizon_url: string;
	    fee: number;
	    sync_frequency: string;
	
	    static createFrom(source: any = {}) {
	        return new StellarConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.network = source["network"];
	        this.horizon_url = source["horizon_url"];
	        this.fee = source["fee"];
	        this.sync_frequency = source["sync_frequency"];
	    }
	}
	export class BlockchainConfig {
	    stellar: StellarConfig;
	    ipfs: IPFSConfig;
	
	    static createFrom(source: any = {}) {
	        return new BlockchainConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stellar = this.convertValues(source["stellar"], StellarConfig);
	        this.ipfs = this.convertValues(source["ipfs"], IPFSConfig);
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
	export class VaultConfig {
	    max_entries: number;
	    auto_sync_enabled: boolean;
	    encryption_scheme: string;
	
	    static createFrom(source: any = {}) {
	        return new VaultConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.max_entries = source["max_entries"];
	        this.auto_sync_enabled = source["auto_sync_enabled"];
	        this.encryption_scheme = source["encryption_scheme"];
	    }
	}
	export class CommitRule {
	    id: number;
	    rule: string;
	    actors: string[];
	
	    static createFrom(source: any = {}) {
	        return new CommitRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.rule = source["rule"];
	        this.actors = source["actors"];
	    }
	}
	export class AppConfig {
	    id: string;
	    repo_id: string;
	    branch: string;
	    tracecore_enabled: boolean;
	    commit_rules: CommitRule[];
	    branching_model: string;
	    encryption_policy: string;
	    actors: string[];
	    federated_providers: string[];
	    default_phase: string;
	    default_vault_path: string;
	    vault_settings: VaultConfig;
	    blockchain: BlockchainConfig;
	    user_id: string;
	    auto_lock_timeout: string;
	    remask_delay: string;
	    theme: string;
	    animations_enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.repo_id = source["repo_id"];
	        this.branch = source["branch"];
	        this.tracecore_enabled = source["tracecore_enabled"];
	        this.commit_rules = this.convertValues(source["commit_rules"], CommitRule);
	        this.branching_model = source["branching_model"];
	        this.encryption_policy = source["encryption_policy"];
	        this.actors = source["actors"];
	        this.federated_providers = source["federated_providers"];
	        this.default_phase = source["default_phase"];
	        this.default_vault_path = source["default_vault_path"];
	        this.vault_settings = this.convertValues(source["vault_settings"], VaultConfig);
	        this.blockchain = this.convertValues(source["blockchain"], BlockchainConfig);
	        this.user_id = source["user_id"];
	        this.auto_lock_timeout = source["auto_lock_timeout"];
	        this.remask_delay = source["remask_delay"];
	        this.theme = source["theme"];
	        this.animations_enabled = source["animations_enabled"];
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
	
	
	
	export class SharingRule {
	    id: number;
	    entry_type: string;
	    targets: string[];
	    encrypted: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SharingRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entry_type = source["entry_type"];
	        this.targets = source["targets"];
	        this.encrypted = source["encrypted"];
	    }
	}
	export class StellarAccountConfig {
	    public_key: string;
	    private_key?: string;
	    EncPassword: number[];
	    EncNonce: number[];
	    EncSalt: number[];
	
	    static createFrom(source: any = {}) {
	        return new StellarAccountConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.public_key = source["public_key"];
	        this.private_key = source["private_key"];
	        this.EncPassword = source["EncPassword"];
	        this.EncNonce = source["EncNonce"];
	        this.EncSalt = source["EncSalt"];
	    }
	}
	
	export class UserConfig {
	    id: string;
	    role: string;
	    signature: string;
	    connected_orgs: string[];
	    stellar_account: StellarAccountConfig;
	    sharing_rules: SharingRule[];
	    two_factor_enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UserConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.role = source["role"];
	        this.signature = source["signature"];
	        this.connected_orgs = source["connected_orgs"];
	        this.stellar_account = this.convertValues(source["stellar_account"], StellarAccountConfig);
	        this.sharing_rules = this.convertValues(source["sharing_rules"], SharingRule);
	        this.two_factor_enabled = source["two_factor_enabled"];
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

export namespace auth {
	
	export class Claims {
	    user_id: string;
	    username: string;
	    email: string;
	    iss?: string;
	    sub?: string;
	    aud?: string[];
	    // Go type: jwt
	    exp?: any;
	    // Go type: jwt
	    nbf?: any;
	    // Go type: jwt
	    iat?: any;
	    jti?: string;
	
	    static createFrom(source: any = {}) {
	        return new Claims(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.username = source["username"];
	        this.email = source["email"];
	        this.iss = source["iss"];
	        this.sub = source["sub"];
	        this.aud = source["aud"];
	        this.exp = this.convertValues(source["exp"], null);
	        this.nbf = this.convertValues(source["nbf"], null);
	        this.iat = this.convertValues(source["iat"], null);
	        this.jti = source["jti"];
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
	export class TokenPairs {
	    access_token: string;
	    refresh_token: string;
	    user_id: string;
	
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
	
	export class CheckEmailResponse {
	    status: string;
	    auth_methods?: string[];
	
	    static createFrom(source: any = {}) {
	        return new CheckEmailResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.auth_methods = source["auth_methods"];
	    }
	}
	export class RecipientPayload {
	    name: string;
	    email: string;
	    role: string;
	
	    static createFrom(source: any = {}) {
	        return new RecipientPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.email = source["email"];
	        this.role = source["role"];
	    }
	}
	export class CreateShareEntryPayload {
	    entry_name: string;
	    entry_type: string;
	    entry_ref: string;
	    status: string;
	    access_mode: string;
	    encryption: string;
	    entry_snapshot: string;
	    expires_at: string;
	    recipients: RecipientPayload[];
	
	    static createFrom(source: any = {}) {
	        return new CreateShareEntryPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry_name = source["entry_name"];
	        this.entry_type = source["entry_type"];
	        this.entry_ref = source["entry_ref"];
	        this.status = source["status"];
	        this.access_mode = source["access_mode"];
	        this.encryption = source["encryption"];
	        this.entry_snapshot = source["entry_snapshot"];
	        this.expires_at = source["expires_at"];
	        this.recipients = this.convertValues(source["recipients"], RecipientPayload);
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
	export class LoginRequest {
	    email?: string;
	    password?: string;
	    publicKey?: string;
	    privateKey?: string;
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
	        this.privateKey = source["privateKey"];
	        this.signedMessage = source["signedMessage"];
	        this.signature = source["signature"];
	    }
	}
	export class LoginResponse {
	    User: models.User;
	    Vault?: models.VaultPayload;
	    Tokens?: auth.TokenPairs;
	    cloud_token: string;
	    vault_runtime_context?: models.VaultRuntimeContext;
	    last_cid: string;
	    dirty: boolean;
	    session_id: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.User = this.convertValues(source["User"], models.User);
	        this.Vault = this.convertValues(source["Vault"], models.VaultPayload);
	        this.Tokens = this.convertValues(source["Tokens"], auth.TokenPairs);
	        this.cloud_token = source["cloud_token"];
	        this.vault_runtime_context = this.convertValues(source["vault_runtime_context"], models.VaultRuntimeContext);
	        this.last_cid = source["last_cid"];
	        this.dirty = source["dirty"];
	        this.session_id = source["session_id"];
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

export namespace main {
	
	export class CheckKeyResponse {
	    id: string;
	    created_at: string;
	    subscription_tier: string;
	    storage_used_gb: number;
	    last_synced_at: string;
	    ok: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CheckKeyResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.created_at = source["created_at"];
	        this.subscription_tier = source["subscription_tier"];
	        this.storage_used_gb = source["storage_used_gb"];
	        this.last_synced_at = source["last_synced_at"];
	        this.ok = source["ok"];
	    }
	}
	export class CreateCheckoutResponse {
	    sessionId: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateCheckoutResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.url = source["url"];
	    }
	}
	export class CreateShareInput {
	    payload: handlers.CreateShareEntryPayload;
	    jwtToken: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateShareInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.payload = this.convertValues(source["payload"], handlers.CreateShareEntryPayload);
	        this.jwtToken = source["jwtToken"];
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
	export class GetSessionResponse {
	    Data: Record<string, any>;
	    Error: any;
	
	    static createFrom(source: any = {}) {
	        return new GetSessionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = source["Data"];
	        this.Error = source["Error"];
	    }
	}
	export class OnboardingStep1Response {
	    identity: string;
	
	    static createFrom(source: any = {}) {
	        return new OnboardingStep1Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.identity = source["identity"];
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
	
	
	
	
	export class User {
	    id: string;
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
	    id: string;
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
	export class Folder {
	    id: string;
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
	export class VaultRuntimeContext {
	    CurrentUser: app_config.UserConfig;
	    AppSettings: app_config.AppConfig;
	    SessionSecrets: Record<string, string>;
	    WorkingBranch: string;
	    LoadedEntries: any[];
	
	    static createFrom(source: any = {}) {
	        return new VaultRuntimeContext(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CurrentUser = this.convertValues(source["CurrentUser"], app_config.UserConfig);
	        this.AppSettings = this.convertValues(source["AppSettings"], app_config.AppConfig);
	        this.SessionSecrets = source["SessionSecrets"];
	        this.WorkingBranch = source["WorkingBranch"];
	        this.LoadedEntries = source["LoadedEntries"];
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

export namespace onboarding_ui_wails {
	
	export class AccountCreationResponse {
	    user_id: string;
	    stellar_key?: string;
	    secret_key?: string;
	
	    static createFrom(source: any = {}) {
	        return new AccountCreationResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.stellar_key = source["stellar_key"];
	        this.secret_key = source["secret_key"];
	    }
	}

}

export namespace onboarding_usecase {
	
	export class AccountCreationRequest {
	    email?: string;
	    password?: string;
	    is_anonymous: boolean;
	    stellar_key?: string;
	
	    static createFrom(source: any = {}) {
	        return new AccountCreationRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.email = source["email"];
	        this.password = source["password"];
	        this.is_anonymous = source["is_anonymous"];
	        this.stellar_key = source["stellar_key"];
	    }
	}
	export class PaymentSetupRequest {
	    user_id: string;
	    tier: string;
	    stripe_payment_method_id?: string;
	    encrypted_payment_data?: string;
	    stellar_public_key?: string;
	    exp?: string;
	    card_number?: string;
	    card_brand: string;
	    exp_month: string;
	    exp_year: string;
	    last_four: string;
	    currency: string;
	    amount: string;
	    payment_method_id: string;
	    cvc?: string;
	    plan: string;
	    product_id: string;
	
	    static createFrom(source: any = {}) {
	        return new PaymentSetupRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.tier = source["tier"];
	        this.stripe_payment_method_id = source["stripe_payment_method_id"];
	        this.encrypted_payment_data = source["encrypted_payment_data"];
	        this.stellar_public_key = source["stellar_public_key"];
	        this.exp = source["exp"];
	        this.card_number = source["card_number"];
	        this.card_brand = source["card_brand"];
	        this.exp_month = source["exp_month"];
	        this.exp_year = source["exp_year"];
	        this.last_four = source["last_four"];
	        this.currency = source["currency"];
	        this.amount = source["amount"];
	        this.payment_method_id = source["payment_method_id"];
	        this.cvc = source["cvc"];
	        this.plan = source["plan"];
	        this.product_id = source["product_id"];
	    }
	}

}

export namespace share_application_use_cases {
	
	export class AddReceiverInput {
	    ShareID: string;
	    Name: string;
	    Email: string;
	    Role: string;
	
	    static createFrom(source: any = {}) {
	        return new AddReceiverInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ShareID = source["ShareID"];
	        this.Name = source["Name"];
	        this.Email = source["Email"];
	        this.Role = source["Role"];
	    }
	}
	export class AddReceiverResult {
	    ShareID: string;
	    RecipientID: string;
	    Message: string;
	
	    static createFrom(source: any = {}) {
	        return new AddReceiverResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ShareID = source["ShareID"];
	        this.RecipientID = source["RecipientID"];
	        this.Message = source["Message"];
	    }
	}
	export class RejectShareResult {
	    ShareID: string;
	    RecipientID: string;
	    Message: string;
	
	    static createFrom(source: any = {}) {
	        return new RejectShareResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ShareID = source["ShareID"];
	        this.RecipientID = source["RecipientID"];
	        this.Message = source["Message"];
	    }
	}

}

export namespace share_domain {
	
	export class EntrySnapshot {
	    entry_name: string;
	    type: string;
	    user_name: string;
	    password: string;
	    website: string;
	    cardholder_name: string;
	    card_number: string;
	    expiry_month: number;
	    expiry_year: number;
	    cvv: string;
	    private_key: string;
	    public_key: string;
	    note: string;
	    genre: string;
	    firstname: string;
	    second_firstname: string;
	    lastname: string;
	    username: string;
	    company: string;
	    social_security_number: string;
	    ID_number: string;
	    driver_license: string;
	    mail: string;
	    telephone: string;
	    address_one: string;
	    address_two: string;
	    address_three: string;
	    city: string;
	    state: string;
	    postal_code: string;
	    country: string;
	    extra_fields: number[];
	
	    static createFrom(source: any = {}) {
	        return new EntrySnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entry_name = source["entry_name"];
	        this.type = source["type"];
	        this.user_name = source["user_name"];
	        this.password = source["password"];
	        this.website = source["website"];
	        this.cardholder_name = source["cardholder_name"];
	        this.card_number = source["card_number"];
	        this.expiry_month = source["expiry_month"];
	        this.expiry_year = source["expiry_year"];
	        this.cvv = source["cvv"];
	        this.private_key = source["private_key"];
	        this.public_key = source["public_key"];
	        this.note = source["note"];
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
	        this.extra_fields = source["extra_fields"];
	    }
	}
	export class Recipient {
	    id: string;
	    share_id: string;
	    name: string;
	    email: string;
	    role: string;
	    // Go type: time
	    joined_at: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    encrypted_blob: number[];
	
	    static createFrom(source: any = {}) {
	        return new Recipient(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.share_id = source["share_id"];
	        this.name = source["name"];
	        this.email = source["email"];
	        this.role = source["role"];
	        this.joined_at = this.convertValues(source["joined_at"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.encrypted_blob = source["encrypted_blob"];
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
	export class ShareEntry {
	    id: string;
	    owner_id: string;
	    entry_name: string;
	    entry_type: string;
	    entry_ref: string;
	    status: string;
	    access_mode: string;
	    encryption: string;
	    entry_snapshot: EntrySnapshot;
	    // Go type: time
	    expires_at?: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    // Go type: time
	    shared_at: any;
	    recipients: Recipient[];
	
	    static createFrom(source: any = {}) {
	        return new ShareEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.owner_id = source["owner_id"];
	        this.entry_name = source["entry_name"];
	        this.entry_type = source["entry_type"];
	        this.entry_ref = source["entry_ref"];
	        this.status = source["status"];
	        this.access_mode = source["access_mode"];
	        this.encryption = source["encryption"];
	        this.entry_snapshot = this.convertValues(source["entry_snapshot"], EntrySnapshot);
	        this.expires_at = this.convertValues(source["expires_at"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.shared_at = this.convertValues(source["shared_at"], null);
	        this.recipients = this.convertValues(source["recipients"], Recipient);
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
	export class ShareAcceptData {
	    share: ShareEntry;
	    recipient: Recipient;
	    blob: number[];
	
	    static createFrom(source: any = {}) {
	        return new ShareAcceptData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.share = this.convertValues(source["share"], ShareEntry);
	        this.recipient = this.convertValues(source["recipient"], Recipient);
	        this.blob = source["blob"];
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

export namespace stellar_recovery_domain {
	
	export class ImportedKey {
	    StellarPublic: string;
	    StellarSecret: string;
	    CanCreate: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ImportedKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.StellarPublic = source["StellarPublic"];
	        this.StellarSecret = source["StellarSecret"];
	        this.CanCreate = source["CanCreate"];
	    }
	}
	export class Subscription {
	    ID: string;
	    UserID: string;
	    Status: string;
	    Tier: string;
	    // Go type: time
	    CreatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Subscription(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.UserID = source["UserID"];
	        this.Status = source["Status"];
	        this.Tier = source["Tier"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
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
	export class Vault {
	    ID: string;
	    // Go type: time
	    CreatedAt: any;
	    StorageUsedGB: number;
	    StorageQuotaGB: number;
	    // Go type: time
	    LastSyncedAt?: any;
	    IPFSNodeID: string;
	    PinataPinID: string;
	
	    static createFrom(source: any = {}) {
	        return new Vault(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.StorageUsedGB = source["StorageUsedGB"];
	        this.StorageQuotaGB = source["StorageQuotaGB"];
	        this.LastSyncedAt = this.convertValues(source["LastSyncedAt"], null);
	        this.IPFSNodeID = source["IPFSNodeID"];
	        this.PinataPinID = source["PinataPinID"];
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
	export class User {
	    ID: string;
	    Email: string;
	    IsAnonymous: boolean;
	    StellarPublicKey: string;
	    EncryptedSecretKey: string;
	    SubscriptionTier: string;
	    SubscriptionID?: string;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Email = source["Email"];
	        this.IsAnonymous = source["IsAnonymous"];
	        this.StellarPublicKey = source["StellarPublicKey"];
	        this.EncryptedSecretKey = source["EncryptedSecretKey"];
	        this.SubscriptionTier = source["SubscriptionTier"];
	        this.SubscriptionID = source["SubscriptionID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
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
	export class RecoveredVault {
	    User?: User;
	    Vault?: Vault;
	    Subscription?: Subscription;
	    SessionToken: string;
	
	    static createFrom(source: any = {}) {
	        return new RecoveredVault(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.User = this.convertValues(source["User"], User);
	        this.Vault = this.convertValues(source["Vault"], Vault);
	        this.Subscription = this.convertValues(source["Subscription"], Subscription);
	        this.SessionToken = source["SessionToken"];
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

export namespace subscription_domain {
	
	export class SubscriptionFeatures {
	    subscription_id: string;
	    StorageGB: number;
	    StorageType: string;
	    CloudBackup: boolean;
	    MobileApps: boolean;
	    SharingLimit: number;
	    UnlimitedSharing: boolean;
	    VersionHistory: boolean;
	    VersionHistoryDays: number;
	    StellarVerification: boolean;
	    Telemetry: boolean;
	    AnonymousAccount: boolean;
	    CryptoPayments: boolean;
	    EncryptedPayments: boolean;
	    Support: string;
	    APIAccess: boolean;
	    Tracecore: boolean;
	    SSO: boolean;
	    TeamFeatures: boolean;
	    payment_method: string;
	    payment_intent: string;
	    browser_extension: boolean;
	    threat_detection: boolean;
	    priority_stellar: boolean;
	    team_size: number;
	    git_cli: boolean;
	    custom_stellar: boolean;
	    on_premise: boolean;
	    multi_cloud: boolean;
	    custom_integrations: boolean;
	    ai_sovereignty: boolean;
	    compliance: string[];
	    sla: string;
	
	    static createFrom(source: any = {}) {
	        return new SubscriptionFeatures(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.subscription_id = source["subscription_id"];
	        this.StorageGB = source["StorageGB"];
	        this.StorageType = source["StorageType"];
	        this.CloudBackup = source["CloudBackup"];
	        this.MobileApps = source["MobileApps"];
	        this.SharingLimit = source["SharingLimit"];
	        this.UnlimitedSharing = source["UnlimitedSharing"];
	        this.VersionHistory = source["VersionHistory"];
	        this.VersionHistoryDays = source["VersionHistoryDays"];
	        this.StellarVerification = source["StellarVerification"];
	        this.Telemetry = source["Telemetry"];
	        this.AnonymousAccount = source["AnonymousAccount"];
	        this.CryptoPayments = source["CryptoPayments"];
	        this.EncryptedPayments = source["EncryptedPayments"];
	        this.Support = source["Support"];
	        this.APIAccess = source["APIAccess"];
	        this.Tracecore = source["Tracecore"];
	        this.SSO = source["SSO"];
	        this.TeamFeatures = source["TeamFeatures"];
	        this.payment_method = source["payment_method"];
	        this.payment_intent = source["payment_intent"];
	        this.browser_extension = source["browser_extension"];
	        this.threat_detection = source["threat_detection"];
	        this.priority_stellar = source["priority_stellar"];
	        this.team_size = source["team_size"];
	        this.git_cli = source["git_cli"];
	        this.custom_stellar = source["custom_stellar"];
	        this.on_premise = source["on_premise"];
	        this.multi_cloud = source["multi_cloud"];
	        this.custom_integrations = source["custom_integrations"];
	        this.ai_sovereignty = source["ai_sovereignty"];
	        this.compliance = source["compliance"];
	        this.sla = source["sla"];
	    }
	}
	export class Subscription {
	    id: string;
	    email: string;
	    wallet?: string;
	    user_id: string;
	    tier: string;
	    expires_at: number;
	    rail: string;
	    tx_hash?: string;
	    active: boolean;
	    activated_at: number;
	    months: number;
	    payment_method: string;
	    payment_intent: string;
	    // Go type: time
	    started_at: any;
	    features: SubscriptionFeatures;
	    ledger: number;
	    billing_cycle: string;
	    trial_ends_at: number;
	    next_billing_date: number;
	    price: number;
	    stripe_subscription_id: string;
	    stellar_payment_address: string;
	    payment_flow_type: string;
	    stellar_schedule_id: string;
	    status: string;
	    ended_at: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Subscription(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.email = source["email"];
	        this.wallet = source["wallet"];
	        this.user_id = source["user_id"];
	        this.tier = source["tier"];
	        this.expires_at = source["expires_at"];
	        this.rail = source["rail"];
	        this.tx_hash = source["tx_hash"];
	        this.active = source["active"];
	        this.activated_at = source["activated_at"];
	        this.months = source["months"];
	        this.payment_method = source["payment_method"];
	        this.payment_intent = source["payment_intent"];
	        this.started_at = this.convertValues(source["started_at"], null);
	        this.features = this.convertValues(source["features"], SubscriptionFeatures);
	        this.ledger = source["ledger"];
	        this.billing_cycle = source["billing_cycle"];
	        this.trial_ends_at = source["trial_ends_at"];
	        this.next_billing_date = source["next_billing_date"];
	        this.price = source["price"];
	        this.stripe_subscription_id = source["stripe_subscription_id"];
	        this.stellar_payment_address = source["stellar_payment_address"];
	        this.payment_flow_type = source["payment_flow_type"];
	        this.stellar_schedule_id = source["stellar_schedule_id"];
	        this.status = source["status"];
	        this.ended_at = source["ended_at"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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

export namespace vaults_domain {
	
	export class Folder {
	    id: string;
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
	export class VaultPayload {
	    version: string;
	    name: string;
	    folders: Folder[];
	    // Go type: Entries
	    entries: any;
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
	        this.entries = this.convertValues(source["entries"], null);
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

