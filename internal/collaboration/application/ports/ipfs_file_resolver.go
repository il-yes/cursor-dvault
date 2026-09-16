package collaboration_ports

import (
	"context"

	vault_dto "vault-app/internal/vault/application/dto"
	vault_queries "vault-app/internal/vault/application/queries"
)

type IPFSFileResolverInterface interface {
	GetIPFSFile(getFilePayload vault_queries.GetIPFSDataQuerry) ([]byte, error)
	GetFileFromIPFS(ctx context.Context, req vault_dto.GetFileFromIPFSRequest) (string, error)
}
