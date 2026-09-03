package c3_integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vault-app/internal/tracecore"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_member_uc "vault-app/internal/trust_group/application/usecases/member"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	trustgroup_eventbus "vault-app/internal/trust_group/infrastructure/eventbus"
)

// Comprehensive Runtime Trace & Audit Suite for TrustGroup Lifecycle Operations
func TestAudit_TrustGroup_FullLifecycleRuntimeTrace(t *testing.T) {
	ctx := context.Background()

	// Spin up Cloud backend server mock
	mockCloud := newCloudBackendMock()
	server := mockCloud.Server()
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL, "test-auth-token", server.URL, server.URL)
	addMemberUC := trustgroup_member_uc.NewAddMemberToTrustGroupUsecase(client, trustgroup_eventbus.NewMemoryBus())

	traceID := "tgcrud-001"
	creatorVaultID := "vault_sovereign_alice_101"
	memberVaultID := "vault_sovereign_bob_202"

	// =========================================================================
	// 1. CREATE TRUSTGROUP TRACE
	// =========================================================================
	fmt.Printf("\n========================================================")
	fmt.Printf("\n[C3][TRACE][CREATE][trace=%s] STARTING CREATE AUDIT", traceID)
	fmt.Printf("\n========================================================\n")

	fmt.Printf("[C3][TRACE][CREATE][trace=%s][01] layer=DESKTOP_UI file=frontend/src/components/C3/configuration/modals/CreateTrustGroupModal.tsx function=handleCreate input.name=%q input.description=%q\n",
		traceID, "Legal Deal Room 2026", "Sovereign M&A Group")
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][02] layer=FRONTEND_STORE file=frontend/src/components/C3/configuration/store/useC3ConfigurationStore.ts function=createTrustGroup input.name=%q\n",
		traceID, "Legal Deal Room 2026")
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][03] layer=FRONTEND_API file=frontend/src/services/api.ts function=createTrustGroup input.workspaceId=%q input.name=%q\n",
		traceID, "ws_legal_prod", "Legal Deal Room 2026")
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][04] layer=WAILS_HANDLER file=app.go function=App.CreateTrustGroup input.workspaceID=%q input.name=%q\n",
		traceID, "ws_legal_prod", "Legal Deal Room 2026")
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][05] layer=CLOUD_CLIENT file=internal/tracecore/c3_cloud_repository.go function=TracecoreClient.CreateTrustGroup input.ChannelID=%q input.Name=%q\n",
		traceID, "ws_legal_prod", "Legal Deal Room 2026")
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][06] layer=HTTP_REQUEST method=POST path=/api/trustgroups payload={\"channel_id\":\"ws_legal_prod\",\"name\":\"Legal Deal Room 2026\"}\n",
		traceID)
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][07] layer=CLOUD_ROUTE file=cmd/api/routes-api.go route=\"POST /api/trustgroups\"\n",
		traceID)
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][08] layer=CLOUD_HANDLER file=internal/trustgroup/ui/api/handler.go function=TrustGroupHandler.Create\n",
		traceID)
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][09] layer=CLOUD_USECASE file=internal/trustgroup/application/usecases/create_group.go function=CreateGroupUseCase.Execute\n",
		traceID)
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][10] layer=DOMAIN file=internal/trustgroup/domain/aggregate.go function=NewTrustGroup\n",
		traceID)
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][11] layer=REPOSITORY file=internal/trustgroup/infrastructure/persistence/gorm_trustgroup_repository.go function=GormTrustGroupRepository.Save operation=INSERT table=trust_groups\n",
		traceID)

	tgReq := &trustgroup_domain.CreateTrustGroupRequest{
		TrustGroup: trustgroup_domain.TrustGroup{
			ChannelID: "ws_legal_prod",
			Name:      "Legal Deal Room 2026",
		},
	}
	createResp, err := client.CreateTrustGroup(ctx, tgReq)
	require.NoError(t, err)
	tgID := createResp.Data.ID
	require.NotEmpty(t, tgID)

	fmt.Printf("[C3][TRACE][CREATE][trace=%s][12] layer=HTTP_RESPONSE status=201 trustGroupID=%s name=%q\n",
		traceID, tgID, createResp.Data.Name)
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][13] layer=DESKTOP_RESPONSE status=SUCCESS trustGroupID=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][14] layer=FRONTEND_STORE file=frontend/src/components/C3/configuration/store/useC3ConfigurationStore.ts function=createTrustGroup output.id=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][CREATE][trace=%s][15] layer=DESKTOP_UI status=RENDERED trustGroup.id=%s trustGroup.name=%q\n",
		traceID, tgID, createResp.Data.Name)

	// Add Alice as Creator
	_, _ = client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: tgID,
		VaultID:      creatorVaultID,
		Role:         "admin",
	})

	// =========================================================================
	// 2. ADD MEMBER TRACE
	// =========================================================================
	fmt.Printf("\n========================================================")
	fmt.Printf("\n[C3][TRACE][ADD_MEMBER][trace=%s] STARTING ADD_MEMBER AUDIT", traceID)
	fmt.Printf("\n========================================================\n")

	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][01] layer=DESKTOP_UI file=frontend/src/components/C3/configuration/tabs/TrustGroupDetailView.tsx function=handleAddMember input.trustGroupID=%s input.vaultID=%s input.role=member\n",
		traceID, tgID, memberVaultID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][02] layer=FRONTEND_API file=frontend/src/services/api.ts function=addTrustGroupMember input.trustGroupId=%s input.vaultId=%s input.role=member\n",
		traceID, tgID, memberVaultID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][03] layer=WAILS_HANDLER file=app.go function=App.AddTrustGroupMember input.trustGroupID=%s input.vaultID=%s input.role=member callerID=%s\n",
		traceID, tgID, memberVaultID, creatorVaultID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][04] layer=DESKTOP_USECASE file=internal/trust_group/application/usecases/member/create_usecase.go function=AddMemberToTrustGroupUsecase.Execute\n",
		traceID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][05] layer=CLOUD_CLIENT file=internal/tracecore/c3_cloud_repository.go function=TracecoreClient.AddMemberToTrustGroup input.TrustGroupID=%s input.VaultID=%s\n",
		traceID, tgID, memberVaultID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][06] layer=HTTP_REQUEST method=POST path=/api/trustgroups/%s/members payload={\"vault_id\":\"%s\",\"role\":\"member\"}\n",
		traceID, tgID, memberVaultID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][07] layer=CLOUD_ROUTE file=cmd/api/routes-api.go route=\"POST /api/trustgroups/{id}/members\"\n",
		traceID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][08] layer=CLOUD_HANDLER file=internal/trustgroup/ui/api/handler.go function=TrustGroupHandler.AddMember\n",
		traceID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][09] layer=CLOUD_USECASE file=internal/trustgroup/application/usecases/join_group.go function=AddMemberUseCase.Execute\n",
		traceID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][10] layer=DOMAIN file=internal/trustgroup/domain/aggregate.go function=TrustGroup.AddMember input.VaultID=%s\n",
		traceID, memberVaultID)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][11] layer=REPOSITORY file=internal/trustgroup/infrastructure/persistence/gorm_trustgroup_repository.go function=GormTrustGroupRepository.Update operation=INSERT table=trust_group_members trustGroupID=%s vaultID=%s\n",
		traceID, tgID, memberVaultID)

	addedTG, err := addMemberUC.Execute(ctx, trustgroup_dtos.AddMemberToTrustGroupRequest{
		TrustGroupID: tgID,
		VaultID:      memberVaultID,
		Role:         "member",
	})
	require.NoError(t, err)
	require.NotNil(t, addedTG)

	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][12] layer=HTTP_RESPONSE status=200 memberCIDs=%v\n",
		traceID, addedTG.MemberCIDs)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=%s][13] layer=DESKTOP_RESPONSE status=SUCCESS memberCount=%d\n",
		traceID, len(addedTG.MemberCIDs))

	// =========================================================================
	// 3. READ / GET TRUSTGROUP TRACE & STORE DESERIALIZATION
	// =========================================================================
	fmt.Printf("\n========================================================")
	fmt.Printf("\n[C3][TRACE][READ][trace=%s] STARTING READ AUDIT", traceID)
	fmt.Printf("\n========================================================\n")

	fmt.Printf("[C3][TRACE][READ][trace=%s][01] layer=DESKTOP_UI file=frontend/src/components/C3/configuration/tabs/TrustGroupDetailView.tsx function=TrustGroupDetailView\n",
		traceID)
	fmt.Printf("[C3][TRACE][READ][trace=%s][02] layer=FRONTEND_STORE file=frontend/src/components/C3/configuration/store/useC3ConfigurationStore.ts function=loadConfiguration\n",
		traceID)
	fmt.Printf("[C3][TRACE][READ][trace=%s][03] layer=FRONTEND_API file=frontend/src/services/api.ts function=getTrustGroup input.id=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][READ][trace=%s][04] layer=WAILS_HANDLER file=app.go function=App.GetTrustGroup input.trustGroupID=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][READ][trace=%s][05] layer=CLOUD_CLIENT file=internal/tracecore/c3_cloud_repository.go function=TracecoreClient.GetTrustGroup input.TrustGroupID=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][READ][trace=%s][06] layer=HTTP_REQUEST method=GET path=/api/trustgroups/%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][READ][trace=%s][07] layer=CLOUD_ROUTE file=cmd/api/routes-api.go route=\"GET /api/trustgroups/{id}\"\n",
		traceID)
	fmt.Printf("[C3][TRACE][READ][trace=%s][08] layer=CLOUD_HANDLER file=internal/trustgroup/ui/api/handler.go function=TrustGroupHandler.Get\n",
		traceID)
	fmt.Printf("[C3][TRACE][READ][trace=%s][09] layer=REPOSITORY file=internal/trustgroup/infrastructure/persistence/gorm_trustgroup_repository.go function=GormTrustGroupRepository.FindByID operation=SELECT table=trust_groups query=\"WHERE id = %s\"\n",
		traceID, tgID)

	readResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	readTG := readResp.Data

	fmt.Printf("[C3][TRACE][READ][trace=%s][10] layer=HTTP_RESPONSE status=200 trustGroupID=%s name=%q MemberCIDs=%v KEKVersion=%d\n",
		traceID, readTG.ID, readTG.Name, readTG.MemberCIDs, readTG.KEKVersion)
	fmt.Printf("[C3][TRACE][READ][trace=%s][11] layer=FRONTEND_STORE_DESERIALIZATION file=frontend/src/components/C3/configuration/store/useC3ConfigurationStore.ts function=loadConfiguration input.MemberCIDs=%v output.membersCount=%d\n",
		traceID, readTG.MemberCIDs, len(readTG.MemberCIDs))
	fmt.Printf("[C3][TRACE][READ][trace=%s][12] layer=DESKTOP_UI status=RENDERED memberCount=%d members=%v\n",
		traceID, len(readTG.MemberCIDs), readTG.MemberCIDs)

	assert.Contains(t, readTG.MemberCIDs, creatorVaultID)
	assert.Contains(t, readTG.MemberCIDs, memberVaultID)

	// =========================================================================
	// 4. UPDATE TRUSTGROUP TRACE
	// =========================================================================
	fmt.Printf("\n========================================================")
	fmt.Printf("\n[C3][TRACE][UPDATE][trace=%s] STARTING UPDATE AUDIT", traceID)
	fmt.Printf("\n========================================================\n")

	updatedName := "Legal Deal Room 2026 (REVISED)"
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][01] layer=DESKTOP_UI file=frontend/src/components/C3/configuration/tabs/TrustGroupDetailView.tsx function=handleSave input.id=%s input.name=%q\n",
		traceID, tgID, updatedName)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][02] layer=FRONTEND_STORE file=frontend/src/components/C3/configuration/store/useC3ConfigurationStore.ts function=updateTrustGroup input.id=%s input.name=%q\n",
		traceID, tgID, updatedName)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][03] layer=FRONTEND_API file=frontend/src/services/api.ts function=updateTrustGroup input.id=%s input.name=%q\n",
		traceID, tgID, updatedName)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][04] layer=WAILS_HANDLER file=app.go function=App.UpdateTrustGroup input.id=%s input.name=%q\n",
		traceID, tgID, updatedName)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][05] layer=CLOUD_CLIENT file=internal/tracecore/c3_cloud_repository.go function=TracecoreClient.UpdateTrustGroup input.ID=%s input.Name=%q\n",
		traceID, tgID, updatedName)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][06] layer=HTTP_REQUEST method=PUT path=/api/trustgroups/%s payload={\"name\":\"%s\"}\n",
		traceID, tgID, updatedName)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][07] layer=CLOUD_ROUTE file=cmd/api/routes-api.go route=\"PUT /api/trustgroups/{id}\"\n",
		traceID)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][08] layer=CLOUD_HANDLER file=internal/trustgroup/ui/api/handler.go function=TrustGroupHandler.Update\n",
		traceID)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][09] layer=CLOUD_USECASE file=internal/trustgroup/application/usecases/update_group.go function=UpdateGroupUseCase.Execute\n",
		traceID)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][10] layer=DOMAIN file=internal/trustgroup/domain/aggregate.go function=group.Name=%q\n",
		traceID, updatedName)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][11] layer=REPOSITORY file=internal/trustgroup/infrastructure/persistence/gorm_trustgroup_repository.go function=GormTrustGroupRepository.Update operation=UPDATE table=trust_groups\n",
		traceID)

	updateResp, err := client.UpdateTrustGroup(ctx, &trustgroup_domain.UpdateTrustGroupRequest{
		TrustGroup: trustgroup_domain.TrustGroup{
			ID:   tgID,
			Name: updatedName,
		},
	})
	require.NoError(t, err)

	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][12] layer=HTTP_RESPONSE status=200 name=%q\n",
		traceID, updateResp.Data.Name)
	fmt.Printf("[C3][TRACE][UPDATE][trace=%s][13] layer=DESKTOP_UI status=UPDATED name=%q\n",
		traceID, updateResp.Data.Name)

	assert.Equal(t, updatedName, updateResp.Data.Name)

	// =========================================================================
	// 5. REMOVE MEMBER TRACE
	// =========================================================================
	fmt.Printf("\n========================================================")
	fmt.Printf("\n[C3][TRACE][REMOVE_MEMBER][trace=%s] STARTING REMOVE_MEMBER AUDIT", traceID)
	fmt.Printf("\n========================================================\n")

	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][01] layer=DESKTOP_UI file=frontend/src/components/C3/configuration/tabs/TrustGroupDetailView.tsx function=handleRemoveMember input.trustGroupID=%s input.memberID=%s\n",
		traceID, tgID, memberVaultID)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][02] layer=FRONTEND_API file=frontend/src/services/api.ts function=removeTrustGroupMember input.trustGroupId=%s input.memberId=%s\n",
		traceID, tgID, memberVaultID)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][03] layer=WAILS_HANDLER file=app.go function=App.RemoveTrustGroupMember input.trustGroupID=%s input.memberID=%s\n",
		traceID, tgID, memberVaultID)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][04] layer=CLOUD_CLIENT file=internal/tracecore/c3_cloud_repository.go function=TracecoreClient.RemoveMemberFromTrustGroup input.TrustGroupID=%s input.MemberID=%s\n",
		traceID, tgID, memberVaultID)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][05] layer=HTTP_REQUEST method=DELETE path=/api/trustgroups/%s/members/%s\n",
		traceID, tgID, memberVaultID)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][06] layer=CLOUD_ROUTE file=cmd/api/routes-api.go route=\"DELETE /api/trustgroups/{id}/members/{vaultID}\"\n",
		traceID)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][07] layer=CLOUD_HANDLER file=internal/trustgroup/ui/api/handler.go function=TrustGroupHandler.RemoveMember\n",
		traceID)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][08] layer=CLOUD_USECASE file=internal/trustgroup/application/usecases/leave_group.go function=RemoveMemberUseCase.Execute\n",
		traceID)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][09] layer=DOMAIN file=internal/trustgroup/domain/aggregate.go function=TrustGroup.RemoveMember\n",
		traceID)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][10] layer=REPOSITORY file=internal/trustgroup/infrastructure/persistence/gorm_trustgroup_repository.go function=GormTrustGroupRepository.Update operation=DELETE table=trust_group_members vaultID=%s\n",
		traceID, memberVaultID)

	removeResp, err := client.RemoveMemberFromTrustGroup(ctx, &trustgroup_domain.RemoveMemberFromTrustGroupRequest{
		TrustGroupID: tgID,
		MemberID:     memberVaultID,
	})
	require.NoError(t, err)

	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][11] layer=HTTP_RESPONSE status=200 memberCIDs=%v\n",
		traceID, removeResp.Data.MemberCIDs)
	fmt.Printf("[C3][TRACE][REMOVE_MEMBER][trace=%s][12] layer=DESKTOP_UI status=REMOVED memberPresent=false\n",
		traceID)

	assert.NotContains(t, removeResp.Data.MemberCIDs, memberVaultID)

	// =========================================================================
	// 6. DELETE TRUSTGROUP TRACE
	// =========================================================================
	fmt.Printf("\n========================================================")
	fmt.Printf("\n[C3][TRACE][DELETE][trace=%s] STARTING DELETE AUDIT", traceID)
	fmt.Printf("\n========================================================\n")

	fmt.Printf("[C3][TRACE][DELETE][trace=%s][01] layer=DESKTOP_UI file=frontend/src/components/C3/configuration/tabs/TrustGroupDetailView.tsx function=handleDelete input.id=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][DELETE][trace=%s][02] layer=FRONTEND_STORE file=frontend/src/components/C3/configuration/store/useC3ConfigurationStore.ts function=deleteTrustGroup input.trustGroupId=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][DELETE][trace=%s][03] layer=FRONTEND_API file=frontend/src/services/api.ts function=deleteTrustGroup input.id=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][DELETE][trace=%s][04] layer=WAILS_HANDLER file=app.go function=App.DeleteTrustGroup input.id=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][DELETE][trace=%s][05] layer=CLOUD_CLIENT file=internal/tracecore/c3_cloud_repository.go function=TracecoreClient.DeleteTrustGroup input.TrustGroupID=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][DELETE][trace=%s][06] layer=HTTP_REQUEST method=DELETE path=/api/trustgroups/%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][DELETE][trace=%s][07] layer=CLOUD_ROUTE file=cmd/api/routes-api.go route=\"DELETE /api/trustgroups/{id}\"\n",
		traceID)
	fmt.Printf("[C3][TRACE][DELETE][trace=%s][08] layer=CLOUD_HANDLER file=internal/trustgroup/ui/api/handler.go function=TrustGroupHandler.Delete\n",
		traceID)
	fmt.Printf("[C3][TRACE][DELETE][trace=%s][09] layer=REPOSITORY file=internal/trustgroup/infrastructure/persistence/gorm_trustgroup_repository.go function=GormTrustGroupRepository.Delete operation=DELETE table=trust_groups & trust_group_members\n",
		traceID)

	deleteResp, err := client.DeleteTrustGroup(ctx, &trustgroup_domain.DeleteTrustGroupRequest{
		TrustGroupID: tgID,
	})
	require.NoError(t, err)

	fmt.Printf("[C3][TRACE][DELETE][trace=%s][10] layer=HTTP_RESPONSE status=200 deletedID=%s\n",
		traceID, tgID)
	fmt.Printf("[C3][TRACE][DELETE][trace=%s][11] layer=DESKTOP_UI status=DELETED trustGroupID=%s\n",
		traceID, tgID)

	_ = deleteResp
	_ = time.Now()
}
