package share_infrastructure

import (
	"context"
	"encoding/json"
	"time"
	share_domain "vault-app/internal/domain/shared"

	"gorm.io/gorm"
)

type GormShareRepository struct {
	db *gorm.DB
}

func NewGormShareRepository(db *gorm.DB) *GormShareRepository {
	return &GormShareRepository{db: db}
}

// ------------------------------
// GORM models (internal l	ayer)
// ------------------------------
type ShareEntryModel struct {
	ID            string `gorm:"primaryKey"`
	OwnerID       string
	EntryRef      string
	EntryName     string
	EntryType     string
	Status        string
	AccessMode    string
	Encryption    string
	EntrySnapshot []byte `gorm:"column:entry_snapshot"` // <---- JSON in DB
	ExpiresAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	SharedAt      time.Time

	Recipients []RecipientModel `gorm:"foreignKey:ShareID;constraint:OnDelete:CASCADE"`
}

func (m *ShareEntryModel) ToDomain() share_domain.ShareEntry {
	var snap share_domain.EntrySnapshot
	_ = json.Unmarshal(m.EntrySnapshot, &snap)

	recipients := make([]share_domain.Recipient, len(m.Recipients))
	for i, r := range m.Recipients {
		recipients[i] = share_domain.Recipient{
			ID:            string(r.ID),
			ShareID:       string(r.ShareID),
			Name:          r.Name,
			Email:         r.Email,
			Role:          r.Role,
			JoinedAt:      r.JoinedAt,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
			EncryptedBlob: r.EncryptedBlob,
		}
	}

	return share_domain.ShareEntry{
		ID:            m.ID,
		OwnerID:       m.OwnerID,
		EntryName:     m.EntryName,
		EntryType:     m.EntryType,
		Status:        m.Status,
		AccessMode:    m.AccessMode,
		Encryption:    m.Encryption,
		EntrySnapshot: snap,
		ExpiresAt:     m.ExpiresAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		SharedAt:      m.SharedAt,
		Recipients:    recipients,
	}
}

func (ShareEntryModel) TableName() string { return "share_entries" }

type RecipientModel struct {
	ID        string `gorm:"primaryKey"`
	ShareID   string
	Name      string
	Email     string
	Role      string
	JoinedAt  time.Time `json:"joined_at" gorm:"column:joined_at"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// Blob containing encrypted vault snapshot (optional)
	EncryptedBlob []byte `gorm:"column:encrypted_blob"`
}

func (m *RecipientModel) ToDomain() share_domain.Recipient {
	return share_domain.Recipient{
		ID:            string(m.ID),
		ShareID:       string(m.ShareID),
		Name:          m.Name,
		Email:         m.Email,
		Role:          m.Role,
		JoinedAt:      m.JoinedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		EncryptedBlob: m.EncryptedBlob,
	}
}

func (RecipientModel) TableName() string { return "recipients" }

// ------------------------------
// Repository implementation
// ------------------------------

func (r *GormShareRepository) ListByUser(userID string) ([]share_domain.ShareEntry, error) {
	var models []ShareEntryModel
	if err := r.db.Preload("Recipients").
		Where("owner_id = ?", userID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return toDomainList(models), nil
}

// ----------------------------------------------
// Returns shares where the user is a recipient
// ----------------------------------------------
func (r *GormShareRepository) ListReceivedByUser(recipientID string) ([]share_domain.ShareEntry, error) {
	var shares []share_domain.ShareEntry

	err := r.db.
		Model(&ShareEntryModel{}).
		Preload("Recipients").
		Where("id IN (SELECT share_id FROM share_recipients WHERE user_id = ?)", recipientID).
		Find(&shares).Error

	if err != nil {
		return nil, err
	}

	return shares, nil
}

func (r *GormShareRepository) GetShareForAccept(
	shareID string,
	recipientUserID string,
) (*share_domain.ShareEntry, *share_domain.Recipient, []byte, error) {

	var share ShareEntryModel
	if err := r.db.
		Preload("Recipients").
		Preload("AuditLog").
		First(&share, "id = ?", shareID).Error; err != nil {
		return nil, nil, nil, err
	}

	var recipient RecipientModel
	if err := r.db.
		First(&recipient,
			"share_id = ? AND user_id = ?",
			shareID, recipientUserID,
		).Error; err != nil {
		return nil, nil, nil, err
	}

	// blob stored directly on recipient row
	blob := recipient.EncryptedBlob

	shareDomain := share.ToDomain()
	recipientDomain := recipient.ToDomain()

	return &shareDomain, &recipientDomain, blob, nil
}

// ---------------------------------------------------------
// Load share + recipient (with encrypted blob)
// ---------------------------------------------------------
func (r *GormShareRepository) GetShareAndRecipient(ctx context.Context, shareID string, userID string) (*share_domain.ShareEntry, *share_domain.Recipient, error) {

	var share share_domain.ShareEntry
	if err := r.db.
		Preload("Recipients").
		Where("id = ?", shareID).
		First(&share).Error; err != nil {
		return nil, nil, share_domain.ErrShareNotFound
	}

	// Find matching recipient
	for _, rcpt := range share.Recipients {
		if rcpt.ID == userID {
			return &share, &rcpt, nil
		}
	}

	return nil, nil, share_domain.ErrRecipientNotAllowed
}


// ---------------------------------------------------------
// Mark a recipient as "accepted"
// ---------------------------------------------------------
func (r *GormShareRepository) MarkRecipientAccepted(ctx context.Context, recipientID string) error {
	return r.db.Exec(`
        UPDATE recipients
        SET role = role, updated_at = NOW()
        WHERE id = ?
    `, recipientID).Error
}
func (r *GormShareRepository) MarkRecipientRejected(ctx context.Context, recipientID string) error {
	return r.db.Exec(`
        UPDATE recipients
        SET role = role, updated_at = NOW(), encrypted_blob = NULL
        WHERE id = ?
    `, recipientID).Error
}

// ---------------------------------------------------------
// Add Recipient
// ---------------------------------------------------------
func (r *GormShareRepository) GetShareByID(ctx context.Context, shareID string) (*share_domain.ShareEntry, error) {
	var share share_domain.ShareEntry
	if err := r.db.WithContext(ctx).
		Preload("Recipients").
		First(&share, shareID).Error; err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *GormShareRepository) CreateRecipient(ctx context.Context, rec *share_domain.Recipient) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

// ---------------------------------------------------------
// Create Share
// ---------------------------------------------------------
// func (r *GormShareRepository) CreateShare(ctx context.Context, s *share_domain.ShareEntry) error {
// 	// Convert domain model to GORM model
// 	model, err := toModel(*s)
// 	if err != nil {
// 		return err
// 	}

// 	// Create in database
// 	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
// 		return err
// 	}

// 	// Update the domain model with the generated ID
// 	s.ID = model.ID

// 	// Update recipient IDs
// 	for i := range s.Recipients {
// 		if i < len(model.Recipients) {
// 			s.Recipients[i].ID = string(model.Recipients[i].ID)
// 		}
// 	}

// 	return nil
// }

func toDomainList(models []ShareEntryModel) []share_domain.ShareEntry {
	res := make([]share_domain.ShareEntry, 0, len(models))
	for _, m := range models {
		res = append(res, m.ToDomain())
	}
	return res
}

func (r *GormShareRepository) GetByID(id string) (*share_domain.ShareEntry, error) {
	var model ShareEntryModel
	if err := r.db.Preload("Recipients").
		First(&model, id).Error; err != nil {
		return nil, err
	}

	domain := model.ToDomain()
	return &domain, nil
}

func (r *GormShareRepository) Save(entry *share_domain.ShareEntry) error {
	model, err := toModel(*entry)
	if err != nil {
		return err
	}
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&model).Error
}

func (r *GormShareRepository) Delete(id string) error {
	return r.db.Delete(&ShareEntryModel{}, id).Error
}

func toModel(e share_domain.ShareEntry) (*ShareEntryModel, error) {
	snapBytes, err := json.Marshal(e.EntrySnapshot)
	if err != nil {
		return nil, err
	}

	recipients := make([]RecipientModel, len(e.Recipients))
	for i, r := range e.Recipients {
		recipients[i] = RecipientModel{
			ID:            r.ID,
			ShareID:       string(r.ShareID),
			Name:          r.Name,
			Email:         r.Email,
			Role:          r.Role,
			JoinedAt:      r.JoinedAt,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
			EncryptedBlob: r.EncryptedBlob,
		}
	}

	return &ShareEntryModel{
		ID:            string(e.ID),
		OwnerID:       e.OwnerID,
		EntryName:     e.EntryName,
		EntryType:     e.EntryType,
		Status:        e.Status,
		AccessMode:    e.AccessMode,
		Encryption:    e.Encryption,
		EntrySnapshot: snapBytes,
		ExpiresAt:     e.ExpiresAt,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
		SharedAt:      e.SharedAt,

		Recipients: recipients,
	}, nil
}
