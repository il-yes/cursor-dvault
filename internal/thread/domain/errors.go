package thread_domain

import "errors"

var (
	ErrRepositoryNil = errors.New("repository is nil")
	ErrThreadNotFound = errors.New("thread not found")
	ErrThreadIDRequired = errors.New("thread id is required")
	ErrThreadTitleRequired = errors.New("thread title is required")
	ErrThreadBusRequired = errors.New("Event bus is nil")
	ErrRequestRequired = errors.New("Request is nil")
	ErrRepositoryResponse = errors.New("thread repository returned nil response")
	ErrThreadRepositoryRequired = errors.New("thread repository is required")
	ErrChannelIDRequired = errors.New("channel id is required")
	ErrAssetTypeRequired = errors.New("asset type is required")
	ErrThreadSubtitleRequired = errors.New("thread subtitle is required")
)