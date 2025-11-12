package api

type SanityApi struct {
	*ApiBase
}

func NewSanityApi(baseURL string) *SanityApi {
	return &SanityApi{
		ApiBase: NewApiBase(baseURL),
	}
}

func (s *SanityApi) CreateNFTMetadata(metadata map[string]interface{}) ([]byte, error) {
	return s.Post("/nft-metadata", metadata)
}

func (s *SanityApi) UpdateNFTMetadata(id string, metadata map[string]interface{}) ([]byte, error) {
	path := "/nft-metadata/" + id
	// You could implement a PUT method in ApiBase or reuse POST logic here
	return s.Post(path, metadata)
}
