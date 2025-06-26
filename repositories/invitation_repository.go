package repositories

import (
	"context"
	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"
	"errors"
	"gorm.io/gorm"
)

type IInvitationRepository interface {
	GetAllInvitations(params queryparams.ListParams) ([]models.Invitation, int64, error)
	GetInvitationByID(id uint) (*models.Invitation, error)
	GetByInvitationKey(ctx context.Context, key string) (*models.Invitation, error) // YENİ METOT
	CreateInvitationWithRelations(ctx context.Context, invitation *models.Invitation) error
	UpdateInvitationWithRelations(ctx context.Context, invitation *models.Invitation) error
	DeleteInvitationWithRelations(ctx context.Context, id uint) error
	GetInvitationCount() (int64, error)
	KeyExists(ctx context.Context, key string) (bool, error)
}

type InvitationRepository struct {
	base IBaseRepository[models.Invitation]
	db   *gorm.DB
}

func NewInvitationRepository() IInvitationRepository {
	db := databaseconfig.GetDB()
	base := NewBaseRepository[models.Invitation](db)
	base.SetAllowedSortColumns([]string{"id", "title", "type", "date", "created_at"})
	base.SetPreloads("InvitationDetail")
	return &InvitationRepository{
		base: base,
		db:   db,
	}
}

func (r *InvitationRepository) GetAllInvitations(params queryparams.ListParams) ([]models.Invitation, int64, error) {
	return r.base.GetAll(params)
}

func (r *InvitationRepository) GetInvitationByID(id uint) (*models.Invitation, error) {
	return r.base.GetByID(id)
}

func (r *InvitationRepository) GetByInvitationKey(ctx context.Context, key string) (*models.Invitation, error) {
	var result models.Invitation
	query := r.db.WithContext(ctx)
	for _, preload := range r.base.(*BaseRepository[models.Invitation]).preloads {
		query = query.Preload(preload)
	}

	err := query.Where("invitation_key = ?", key).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound 
		}
		return nil, err
	}
	return &result, nil
}

func (r *InvitationRepository) CreateInvitationWithRelations(ctx context.Context, invitation *models.Invitation) error {
	return r.base.CreateWithRelations(ctx, invitation)
}

func (r *InvitationRepository) UpdateInvitationWithRelations(ctx context.Context, invitation *models.Invitation) error {
	return r.base.UpdateWithRelations(ctx, invitation)
}

func (r *InvitationRepository) DeleteInvitationWithRelations(ctx context.Context, id uint) error {
	return r.base.DeleteWithRelations(ctx, id)
}

func (r *InvitationRepository) GetInvitationCount() (int64, error) {
	return r.base.GetCount()
}

func (r *InvitationRepository) KeyExists(ctx context.Context, key string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Invitation{}).Where("invitation_key = ?", key).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}