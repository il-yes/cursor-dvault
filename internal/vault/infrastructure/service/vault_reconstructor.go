package vaults_service

import (
	"context"
	"encoding/json"
	"fmt"

	"vault-app/internal/utils"
	vault_queries "vault-app/internal/vault/application/queries"
	vaults_domain "vault-app/internal/vault/domain"
)

type QueryExecutor interface {
	Execute(ctx context.Context, cmd vault_queries.GetIPFSDataQuerry) (*vault_queries.GetIPFSDataResponse, error)
}

type VaultReconstructor struct {
	Query QueryExecutor
}

func NewVaultReconstructor(q *vault_queries.GetIPFSDataQuerryHandler) *VaultReconstructor {
	return &VaultReconstructor{
		Query: q,
	}
}

func (r *VaultReconstructor) BuildFromRoot(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
) (vaults_domain.VaultPayload, error) {
	utils.LogPretty("VaultReconstructor - BuildFromRoot -  cmd", "ok")

	// -----------------------------
	// 1. ROOT
	// -----------------------------
	rootRes, err := r.Query.Execute(ctx, cmd)
	if err != nil {
		utils.LogPretty("VaultReconstructor - BuildFromRoot - fail get ipfs query", err)
		return vaults_domain.VaultPayload{}, err
	}
	// utils.LogPretty("VaultReconstructor - BuildFromRoot - rootRes", rootRes)

	root := rootRes.NodeBeta
	utils.LogPretty("VaultReconstructor - BuildFromRoot - root", root)
	cmd.CID = root.Personal.CID
	personalRes, err := r.Query.Execute(ctx, cmd)
	if err != nil {
		utils.LogPretty("VaultReconstructor - BuildFromRoot - fail get ipfs query", err)
		return vaults_domain.VaultPayload{}, err
	}
	var personalRoot vaults_domain.PersonalNode
	if err := json.Unmarshal(personalRes.Raw, &personalRoot); err != nil {
		utils.LogPretty("VaultReconstructor - BuildFromRoot - fail to  convert ipfs data to VaultNode", err)
		return vaults_domain.VaultPayload{}, err
	}

	cmd.CID = root.Collaborative.CID
	collaborativeRes, err := r.Query.Execute(ctx, cmd)
	if err != nil {
		utils.LogPretty("VaultReconstructor - BuildFromRoot - fail get ipfs query", err)
		return vaults_domain.VaultPayload{}, err
	}
	var collaborativeRoot vaults_domain.CollaborativeNode
	if err := json.Unmarshal(collaborativeRes.Raw, &collaborativeRoot); err != nil {
		utils.LogPretty("VaultReconstructor - BuildFromRoot - fail to  convert ipfs data to CollaborativeNode", err)
		return vaults_domain.VaultPayload{}, err
	}
	utils.LogPretty( "collaborativeRoot after unmarshal", collaborativeRoot)

	// TODO: change the condition
	if root.Version == "" {
		// FIRST SYNC CASE
		emptyVp := emptyVaultPayload(cmd.VaultName, "1.0.0")
		return vaults_domain.VaultPayload{
			Version: rootRes.Data.Version,
			Name:    rootRes.Data.Name,
			// BaseVaultContent: vaults_domain.BaseVaultContent{
			// 	Entries: emptyVp.Entries,
			// 	Folders: emptyVp.Folders, // TODO: risky...
			// },
			Personal: vaults_domain.BaseVaultContent{
				Entries: emptyVp.Entries,
				Folders: emptyVp.Folders, // TODO: risky...
			},
			Collaborative: vaults_domain.C3VaultContent{
				Workspaces:        emptyVp.Collaborative.Workspaces,
				Channels:          emptyVp.Collaborative.Channels,
				Threads:           emptyVp.Collaborative.Threads,
				ShareEntries:      emptyVp.Collaborative.ShareEntries,
				TrustGroups:       emptyVp.Collaborative.TrustGroups,
				TrustGroupMembers: emptyVp.Collaborative.TrustGroupMembers,
				Federation:        emptyVp.Collaborative.Federation,
				Participants:      emptyVp.Collaborative.Participants,
				Assets:            emptyVp.Collaborative.Assets,
				Index:             emptyVp.Collaborative.Index,
			},
		}, nil
	}

	personalVault, err := r.resolvePersonalPart(ctx, cmd, personalRoot)
	if err != nil {
		utils.LogPretty("VaultReconstructor - personalVault - failed to reconstruct personalVault %v", err)
		return vaults_domain.VaultPayload{}, err
	}
	collaborativeVault, err := r.resolveCollaborativePart(ctx, cmd, collaborativeRoot)
	if err != nil {
		utils.LogPretty("VaultReconstructor - collaborativeVault - failed to reconstruct collaborativeVault", err)
		return vaults_domain.VaultPayload{}, err
	}

	return vaults_domain.VaultPayload{
		Version: root.Version,
		Name:    "reconstructed",

		// BaseVaultContent: vaults_domain.BaseVaultContent{
		// 	Entries:     entries,
		// 	Folders:     folders,
		// 	Attachments: attachments,
		// 	Index:       index,
		// },
		Personal: vaults_domain.BaseVaultContent{
			Entries:     personalVault.entries,
			Folders:     personalVault.folders,
			Attachments: personalVault.attachments,
			Index:       personalVault.index,
		},
		Collaborative: vaults_domain.C3VaultContent{
			Workspaces:        collaborativeVault.Workspaces,
			Channels:          collaborativeVault.Channels,
			Threads:           collaborativeVault.Threads,
			ShareEntries:      collaborativeVault.ShareEntries,
			TrustGroups:       collaborativeVault.TrustGroups,
			TrustGroupMembers: collaborativeVault.TrustGroupMembers,
			Federation:        collaborativeVault.Federation,
			Participants:      collaborativeVault.Participants,
			Assets:            collaborativeVault.Assets,
			Index:             collaborativeVault.Index,
		},
	}, nil
}

type ResolvePersonalResponse struct {
	entries     vaults_domain.Entries
	folders     []vaults_domain.Folder
	attachments []vaults_domain.Attachment
	index       vaults_domain.Index
}

func (r *VaultReconstructor) resolvePersonalPart(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	root vaults_domain.PersonalNode,
) (*ResolvePersonalResponse, error) {
	// -----------------------------
	// 2. ENTRIES ROOT
	// -----------------------------
	entriesRes, err := r.Query.Execute(ctx, cmd.WithCID(root.Entries.CID))
	if err != nil {
		return nil, err
	}

	var entriesRoot vaults_domain.EntriesRoot
	if err := json.Unmarshal(entriesRes.Raw, &entriesRoot); err != nil {
		return nil, err
	}

	// -----------------------------
	// 3. FOLDERS ROOT
	// -----------------------------
	foldersRes, err := r.Query.Execute(ctx, cmd.WithCID(root.Folders.CID))
	if err != nil {
		return nil, err
	}

	var foldersRoot vaults_domain.FoldersRoot
	if err := json.Unmarshal(foldersRes.Raw, &foldersRoot); err != nil {
		return nil, err
	}

	// -----------------------------
	// 4. ATTACHMENT ROOT
	// -----------------------------
	attachmentsRes, err := r.Query.Execute(ctx, cmd.WithCID(root.Attachments.CID))
	if err != nil {
		return nil, err
	}

	var attachmentsRoot vaults_domain.AttachementsRoot
	if err := json.Unmarshal(attachmentsRes.Raw, &attachmentsRoot); err != nil {
		return nil, err
	}

	// -----------------------------
	// 5. Resolve entries (CORRECT)
	// -----------------------------
	entries, err := r.resolveEntries(ctx, cmd, entriesRoot)
	if err != nil {
		return nil, err
	}

	// -----------------------------
	// 6. Resolve folders (CORRECT)
	// -----------------------------
	folders, err := r.resolveFolders(ctx, cmd, foldersRoot)
	if err != nil {
		return nil, err
	}

	// -----------------------------
	// 7. Resolve attachments (CORRECT)
	// -----------------------------
	attachments, err := r.resolveAttachments(ctx, cmd, attachmentsRoot)
	if err != nil {
		return nil, err
	}

	// -----------------------------
	// 8. Resolve index (CORRECT)
	// -----------------------------
	index, err := r.resolveIndex(ctx, cmd, root.Index)
	if err != nil {
		return nil, err
	}

	return &ResolvePersonalResponse{
		entries:     entries,
		folders:     folders,
		attachments: attachments,
		index:       index,
	}, err
}

type ResolveCollaborativeResponse struct {
	Workspaces        []vaults_domain.Workspace
	Channels          []vaults_domain.Channel
	Threads           []vaults_domain.Thread
	ShareEntries      []vaults_domain.ShareEntry
	TrustGroups       []vaults_domain.TrustGroup
	TrustGroupMembers []vaults_domain.TrustGroupMember
	Federation        vaults_domain.FederationSnapshot
	Participants      []vaults_domain.Participant
	Assets            []vaults_domain.Asset
	Index             vaults_domain.IndexC3
}

func (r *VaultReconstructor) resolveCollaborativePart(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	root vaults_domain.CollaborativeNode,
) (*ResolveCollaborativeResponse, error) {
	// -----------------------------
	// 1. WORKSPACE ROOT
	// -----------------------------
	workspacesRes, err := r.Query.Execute(ctx, cmd.WithCID(root.Workspaces.CID))
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  workspacesRes", workspacesRes)
		return nil, err
	}

	var workspacesRoot vaults_domain.WorkspacesRoot
	if err := json.Unmarshal(workspacesRes.Raw, &workspacesRoot); err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  workspacesRoot", workspacesRoot)
		return nil, err
	}
	// -----------------------------
	// 1.2 Resolve WORKSPACE
	// -----------------------------
	workspaces, err := r.resolveWorkspaces(ctx, cmd, workspacesRoot)
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  workspaces", workspaces)
		return nil, err
	}

	// -----------------------------
	// 2. CHANNEL ROOT
	// -----------------------------
	channelsRes, err := r.Query.Execute(ctx, cmd.WithCID(root.Channels.CID))
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  channelsRes", channelsRes)
		return nil, err
	}

	var channelsRoot vaults_domain.ChannelsRoot
	if err := json.Unmarshal(channelsRes.Raw, &channelsRoot); err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  channelsRoot", channelsRoot)
		return nil, err
	}
	// -----------------------------
	// 2.2 Resolve CHANNEL
	// -----------------------------
	channels, err := r.resolveChannels(ctx, cmd, channelsRoot)
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  channels", channels)
		return nil, err
	}

	// -----------------------------
	// 3. THREAD ROOT
	// -----------------------------
	threadsRes, err := r.Query.Execute(ctx, cmd.WithCID(root.Threads.CID))
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  threadsRes", threadsRes)
		return nil, err
	}

	var threadsRoot vaults_domain.ThreadsRoot
	if err := json.Unmarshal(threadsRes.Raw, &threadsRoot); err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  threadsRoot", threadsRoot)
		return nil, err
	}
	// -----------------------------
	// 3.2 Resolve THREAD
	// -----------------------------
	threads, err := r.resolveThreads(ctx, cmd, threadsRoot)
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  threads", threads)
		return nil, err
	}

	// -----------------------------
	// 4. ASSET ROOT
	// -----------------------------
	assetsRes, err := r.Query.Execute(ctx, cmd.WithCID(root.Assets.CID))
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  assetsRes", assetsRes)
		return nil, err
	}

	var assetsRoot vaults_domain.AssetsRoot
	if err := json.Unmarshal(assetsRes.Raw, &assetsRoot); err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  assetsRoot", assetsRoot)
		return nil, err
	}
	// -----------------------------
	// 4.2 Resolve ASSET
	// -----------------------------
	assets, err := r.resolveAssets(ctx, cmd, assetsRoot)
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  assets", assets)
		return nil, err
	}

	// -----------------------------
	// 4. TRUSTGROUP ROOT
	// -----------------------------
	trustgroupsRes, err := r.Query.Execute(ctx, cmd.WithCID(root.TrustGroups.CID))
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  trustgroupsRes", trustgroupsRes)
		return nil, err
	}

	var trustgroupsRoot vaults_domain.TrustGroupsRoot
	if err := json.Unmarshal(trustgroupsRes.Raw, &trustgroupsRoot); err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  trustgroupsRoot", trustgroupsRoot)
		return nil, err
	}
	// -----------------------------
	// 4.2 Resolve TRUSTGROUP
	// -----------------------------
	trustgroups, err := r.resolveTrustGroups(ctx, cmd, trustgroupsRoot)
	if err != nil {
		utils.LogPretty("VaultReconstructor - resolveCollaborativePart - failed to get  trustgroups", trustgroups)
		return nil, err
	}

	return &ResolveCollaborativeResponse{
		Workspaces:        workspaces,
		Channels:          channels,
		Threads:           threads,
		ShareEntries:      []vaults_domain.ShareEntry{},
		TrustGroups:       trustgroups,
		TrustGroupMembers: []vaults_domain.TrustGroupMember{},
		Federation:        vaults_domain.FederationSnapshot{},
		Participants:      []vaults_domain.Participant{},
		Assets:            assets,
		Index:             vaults_domain.IndexC3{},
	}, nil
}

func (r *VaultReconstructor) resolveEntries(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	entriesRoot vaults_domain.EntriesRoot,
) (vaults_domain.Entries, error) {

	var result vaults_domain.Entries

	for _, link := range entriesRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		// 1. Detect type first (light struct)
		var meta struct {
			Type string `json:"type"`
		}

		if err := json.Unmarshal(res.Raw, &meta); err != nil {
			return result, err
		}

		// 2. Dispatch by type (like reverse BuildEntries)
		switch meta.Type {

		case "login":
			var e vaults_domain.LoginEntry
			if err := json.Unmarshal(res.Raw, &e); err != nil {
				return result, err
			}
			result.Login = append(result.Login, e)

		case "card":
			var e vaults_domain.CardEntry
			if err := json.Unmarshal(res.Raw, &e); err != nil {
				return result, err
			}
			result.Card = append(result.Card, e)

		case "identity":
			var e vaults_domain.IdentityEntry
			if err := json.Unmarshal(res.Raw, &e); err != nil {
				return result, err
			}
			result.Identity = append(result.Identity, e)

		case "note":
			var e vaults_domain.NoteEntry
			if err := json.Unmarshal(res.Raw, &e); err != nil {
				return result, err
			}
			result.Note = append(result.Note, e)

		case "sshkey":
			var e vaults_domain.SSHKeyEntry
			if err := json.Unmarshal(res.Raw, &e); err != nil {
				return result, err
			}
			result.SSHKey = append(result.SSHKey, e)

		default:
			return result, fmt.Errorf("unknown entry type: %s", meta.Type)
		}
	}

	return result, nil
}

func (r *VaultReconstructor) resolveFolders(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	foldersRoot vaults_domain.FoldersRoot,
) ([]vaults_domain.Folder, error) {

	var result []vaults_domain.Folder

	for _, link := range foldersRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var folder vaults_domain.Folder
		if err := json.Unmarshal(res.Raw, &folder); err != nil {
			return result, err
		}

		result = append(result, folder)
	}

	return result, nil
}
func (r *VaultReconstructor) resolveAttachments(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	attachmentsRoot vaults_domain.AttachementsRoot,
) ([]vaults_domain.Attachment, error) {

	var result []vaults_domain.Attachment

	for _, link := range attachmentsRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var attachment vaults_domain.Attachment
		if err := json.Unmarshal(res.Raw, &attachment); err != nil {
			return result, err
		}

		result = append(result, attachment)
	}

	return result, nil
}

func (r *VaultReconstructor) resolveIndex(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	indexLink vaults_domain.Link,
) (vaults_domain.Index, error) {

	res, err := r.Query.Execute(ctx, cmd.WithCID(indexLink.CID))
	if err != nil {
		return vaults_domain.Index{}, err
	}

	var index vaults_domain.Index
	if err := json.Unmarshal(res.Raw, &index); err != nil {
		return vaults_domain.Index{}, err
	}

	return index, nil
}

func emptyVaultPayload(name string, version string) vaults_domain.VaultPayload {
	return *vaults_domain.InitEmptyVaultPayload(name, version)
}
