package services

import (
	"context"
	"crypto/rand"
	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"
	"davet.link/repositories"
	"errors"
	"go.uber.org/zap"
)

type IInvitationService interface {
	GetAllInvitations(params queryparams.ListParams) (*queryparams.PaginatedResult, error)
	GetInvitationByID(id uint) (*models.Invitation, error)
	GetInvitationByKey(ctx context.Context, key string) (*models.Invitation, error) // YENİ METOT
	CreateInvitationWithRelations(ctx context.Context, invitation *models.Invitation) error
	UpdateInvitationWithRelations(ctx context.Context, invitation *models.Invitation) error
	DeleteInvitationWithRelations(ctx context.Context, id uint) error
	GetInvitationCount() (int64, error)
}

type InvitationService struct {
	repo repositories.IInvitationRepository
}

func NewInvitationService() IInvitationService {
	return &InvitationService{repo: repositories.NewInvitationRepository()}
}

func (s *InvitationService) GetAllInvitations(params queryparams.ListParams) (*queryparams.PaginatedResult, error) {
	invitations, totalCount, err := s.repo.GetAllInvitations(params)
	if err != nil {
		logconfig.Log.Error("Davetiyeler alınamadı", zap.Error(err))
		return nil, errors.New("davetiyeler getirilirken bir veritabanı hatası oluştu")
	}
	result := &queryparams.PaginatedResult{
		Data: invitations,
		Meta: queryparams.PaginationMeta{CurrentPage: params.Page, PerPage: params.PerPage, TotalItems: totalCount, TotalPages: queryparams.CalculateTotalPages(totalCount, params.PerPage)},
	}
	return result, nil
}

func (s *InvitationService) GetInvitationByID(id uint) (*models.Invitation, error) {
	invitation, err := s.repo.GetInvitationByID(id)
	if err != nil {
		logconfig.Log.Warn("Davetiye ID ile bulunamadı", zap.Uint("id", id), zap.Error(err))
		return nil, errors.New("belirtilen ID ile davetiye bulunamadı")
	}
	return invitation, nil
}

func (s *InvitationService) GetInvitationByKey(ctx context.Context, key string) (*models.Invitation, error) {
	invitation, err := s.repo.GetByInvitationKey(ctx, key)
	if err != nil {
		logconfig.Log.Warn("Davetiye anahtar ile bulunamadı", zap.String("key", key), zap.Error(err))
		return nil, errors.New("belirtilen anahtar ile davetiye bulunamadı")
	}
	return invitation, nil
}

func (s *InvitationService) CreateInvitationWithRelations(ctx context.Context, invitation *models.Invitation) error {
	for {
		key := generateInvitationKey(11)
		exists, err := s.repo.KeyExists(ctx, key)
		if err != nil {
			logconfig.Log.Error("InvitationKey kontrolü sırasında veritabanı hatası", zap.Error(err))
			return errors.New("davetiye anahtarı kontrol edilemedi")
		}
		if !exists {
			invitation.InvitationKey = key
			break
		}
	}
	return s.repo.CreateInvitationWithRelations(ctx, invitation)
}

func (s *InvitationService) UpdateInvitationWithRelations(ctx context.Context, invitation *models.Invitation) error {
	return s.repo.UpdateInvitationWithRelations(ctx, invitation)
}

func (s *InvitationService) DeleteInvitationWithRelations(ctx context.Context, id uint) error {
	return s.repo.DeleteInvitationWithRelations(ctx, id)
}

func (s *InvitationService) GetInvitationCount() (int64, error) {
	return s.repo.GetInvitationCount()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateInvitationKey(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("Kriptografik anahtar üretilemedi: " + err.Error())
	}
	for i, v := range b {
		b[i] = letterBytes[v%byte(len(letterBytes))]
	}
	return string(b)
}
var _ IInvitationService = (*InvitationService)(nil)