package vault_ui




type VaultHandler struct {
	OpenVaultHandler *OpenVaultHandler	
}	

func NewVaultHandler(openVaultHandler *OpenVaultHandler) *VaultHandler {
	return &VaultHandler{
		OpenVaultHandler: openVaultHandler,
	}
}