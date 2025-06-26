package handlers

import (
	"context"
	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/filemanager"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/requests"
	"davet.link/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type DashboardInvitationHandler struct {
	invitationService services.IInvitationService
	categoryService   services.IInvitationCategoryService
}

func NewDashboardInvitationHandler() *DashboardInvitationHandler {
	return &DashboardInvitationHandler{
		invitationService: services.NewInvitationService(),
		categoryService:   services.NewInvitationCategoryService(),
	}
}

func (h *DashboardInvitationHandler) ShowInvitation(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz davetiye anahtarı.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	invitation, err := h.invitationService.GetInvitationByKey(c.UserContext(), key)
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/invitations/show", "layouts/dashboard", fiber.Map{
		"Title":      "Davetiye Detayları",
		"Invitation": invitation,
	})
}

func (h *DashboardInvitationHandler) ListInvitations(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
		logconfig.Log.Warn("Davetiye listesi: Query parametreleri parse edilemedi.", zap.Error(err))
		params = queryparams.DefaultListParams()
	}
	if params.Page <= 0 {
		params.Page = queryparams.DefaultPage
	}
	if params.PerPage <= 0 || params.PerPage > queryparams.MaxPerPage {
		params.PerPage = queryparams.DefaultPerPage
	}
	if params.SortBy == "" {
		params.SortBy = queryparams.DefaultSortBy
	}
	if params.OrderBy == "" {
		params.OrderBy = queryparams.DefaultOrderBy
	}
	result, err := h.invitationService.GetAllInvitations(params)
	renderData := fiber.Map{
		"Title":  "Davetiyeler",
		"Result": result,
		"Params": params,
	}
	if err != nil {
		logconfig.Log.Error("Davetiye listesi DB Hatası", zap.Error(err))
		renderData[renderer.FlashErrorKeyView] = "Davetiyeler getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.Invitation{},
			Meta: queryparams.PaginationMeta{CurrentPage: params.Page, PerPage: params.PerPage},
		}
	}
	return renderer.Render(c, "dashboard/invitations/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *DashboardInvitationHandler) ShowCreateInvitation(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAllCategories(queryparams.DefaultListParams())
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kategoriler getirilemedi.")
	}
	return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
		"Title":      "Yeni Davetiye Oluştur",
		"Categories": categories,
	})
}

func (h *DashboardInvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	if err := requests.ValidateInvitationRequest(c); err != nil {
		req, _ := c.Locals("invitationRequest").(requests.InvitationRequest)
		// DÜZELTME: Fonksiyon doğru parametre ile çağrıldı.
		categories, _ := h.categoryService.GetAllCategories(queryparams.DefaultListParams())
		return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
			"Title": "Yeni Davetiye Oluştur", "Categories": categories, "FormData": req,
		})
	}

	req := c.Locals("invitationRequest").(requests.InvitationRequest)
	userID, _ := c.Locals("userID").(uint)

	newFileName, err := filemanager.UploadFile(c, "image", "invitations")
	if err != nil && err != filemanager.ErrFileNotProvided {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Resim yüklenemedi: "+err.Error())
		// DÜZELTME: Fonksiyon doğru parametre ile çağrıldı.
		categories, _ := h.categoryService.GetAllCategories(queryparams.DefaultListParams())
		return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
			"Title": "Yeni Davetiye Oluştur", "Categories": categories, "FormData": req,
		})
	}

	invitation := &models.Invitation{
		UserID:        userID,
		CategoryID:    req.CategoryID,
		Template:      req.Template,
		Type:          req.Type,
		Title:         req.Title,
		Image:         newFileName,
		Venue:         req.Venue,
		Address:       req.Address,
		Location:      req.Location,
		Telephone:     req.Telephone,
		Date:          req.Date,
		Time:          req.Time,
		IsConfirmed:   req.IsConfirmed,
		IsParticipant: req.IsParticipant,
	}

	detail := &models.InvitationDetail{
		Title:              req.Detail.Title,
		BrideName:          req.Detail.BrideName,
		BrideSurname:       req.Detail.BrideSurname,
		BrideMotherName:    req.Detail.BrideMotherName,
		BrideMotherSurname: req.Detail.BrideMotherSurname,
		BrideFatherName:    req.Detail.BrideFatherName,
		BrideFatherSurname: req.Detail.BrideFatherSurname,
		GroomName:          req.Detail.GroomName,
		GroomSurname:       req.Detail.GroomSurname,
		GroomMotherName:    req.Detail.GroomMotherName,
		GroomMotherSurname: req.Detail.GroomMotherSurname,
		GroomFatherName:    req.Detail.GroomFatherName,
		GroomFatherSurname: req.Detail.GroomFatherSurname,
		Person:             req.Detail.Person,
		MotherName:         req.Detail.MotherName,
		MotherSurname:      req.Detail.MotherSurname,
		FatherName:         req.Detail.FatherName,
		FatherSurname:      req.Detail.FatherSurname,
		IsMotherLive:       req.Detail.IsMotherLive,
		IsFatherLive:       req.Detail.IsFatherLive,
		IsBrideMotherLive:  req.Detail.IsBrideMotherLive,
		IsBrideFatherLive:  req.Detail.IsBrideFatherLive,
		IsGroomMotherLive:  req.Detail.IsGroomMotherLive,
		IsGroomFatherLive:  req.Detail.IsGroomFatherLive,
	}
	invitation.InvitationDetail = detail

	if err := h.invitationService.CreateInvitationWithRelations(c.UserContext(), invitation); err != nil {
		if newFileName != "" {
			// DÜZELTME: Fonksiyonun dönüş değeri kullanılmadı.
			filemanager.DeleteFile("invitations", newFileName)
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye oluşturulamadı: "+err.Error())
		// DÜZELTME: Fonksiyon doğru parametre ile çağrıldı.
		categories, _ := h.categoryService.GetAllCategories(queryparams.DefaultListParams())
		return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
			"Title": "Yeni Davetiye Oluştur", "Categories": categories, "FormData": req,
		})
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}

func (h *DashboardInvitationHandler) ShowUpdateInvitation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz davetiye ID'si.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	invitation, err := h.invitationService.GetInvitationByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	// DÜZELTME: Fonksiyon doğru parametre ile çağrıldı.
	categories, err := h.categoryService.GetAllCategories(queryparams.DefaultListParams())
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kategoriler getirilemedi.")
	}
	return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
		"Title":      "Davetiye Düzenle",
		"FormData":   invitation,
		"Categories": categories,
	})
}

func (h *DashboardInvitationHandler) UpdateInvitation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz davetiye ID'si.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	redirectURL := fmt.Sprintf("/dashboard/invitations/update/%d", id)

	if err := requests.ValidateInvitationRequest(c); err != nil {
		req, _ := c.Locals("invitationRequest").(requests.InvitationRequest)
		// DÜZELTME: Fonksiyon doğru parametre ile çağrıldı.
		categories, _ := h.categoryService.GetAllCategories(queryparams.DefaultListParams())
		return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
			"Title": "Davetiye Düzenle", "Categories": categories, "FormData": req, "ID": id,
		})
	}

	existingInvitation, err := h.invitationService.GetInvitationByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Güncellenecek davetiye bulunamadı.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	req := c.Locals("invitationRequest").(requests.InvitationRequest)
	userID, _ := c.Locals("userID").(uint)

	newFileName, err := filemanager.UploadFile(c, "image", "invitations")
	if err != nil && err != filemanager.ErrFileNotProvided {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yeni resim yüklenemedi: "+err.Error())
		return c.Redirect(redirectURL, http.StatusSeeOther)
	}
	var oldPhotoToDelete string
	if newFileName != "" {
		oldPhotoToDelete = existingInvitation.Image
		existingInvitation.Image = newFileName
	}

	existingInvitation.CategoryID = req.CategoryID
	existingInvitation.Template = req.Template
	existingInvitation.Type = req.Type
	existingInvitation.Title = req.Title
	existingInvitation.Venue = req.Venue
	existingInvitation.Address = req.Address
	existingInvitation.Location = req.Location
	existingInvitation.Telephone = req.Telephone
	existingInvitation.Date = req.Date
	existingInvitation.Time = req.Time
	existingInvitation.IsConfirmed = req.IsConfirmed
	existingInvitation.IsParticipant = req.IsParticipant
	existingInvitation.UpdatedBy = userID

	if existingInvitation.InvitationDetail != nil {
		existingInvitation.InvitationDetail.Title = req.Detail.Title
		existingInvitation.InvitationDetail.BrideName = req.Detail.BrideName
		existingInvitation.InvitationDetail.BrideSurname = req.Detail.BrideSurname
		existingInvitation.InvitationDetail.BrideMotherName = req.Detail.BrideMotherName
		existingInvitation.InvitationDetail.BrideMotherSurname = req.Detail.BrideMotherSurname
		existingInvitation.InvitationDetail.BrideFatherName = req.Detail.BrideFatherName
		existingInvitation.InvitationDetail.BrideFatherSurname = req.Detail.BrideFatherSurname
		existingInvitation.InvitationDetail.GroomName = req.Detail.GroomName
		existingInvitation.InvitationDetail.GroomSurname = req.Detail.GroomSurname
		existingInvitation.InvitationDetail.GroomMotherName = req.Detail.GroomMotherName
		existingInvitation.InvitationDetail.GroomMotherSurname = req.Detail.GroomMotherSurname
		existingInvitation.InvitationDetail.GroomFatherName = req.Detail.GroomFatherName
		existingInvitation.InvitationDetail.GroomFatherSurname = req.Detail.GroomFatherSurname
		existingInvitation.InvitationDetail.Person = req.Detail.Person
		existingInvitation.InvitationDetail.MotherName = req.Detail.MotherName
		existingInvitation.InvitationDetail.MotherSurname = req.Detail.MotherSurname
		existingInvitation.InvitationDetail.FatherName = req.Detail.FatherName
		existingInvitation.InvitationDetail.FatherSurname = req.Detail.FatherSurname
		existingInvitation.InvitationDetail.IsMotherLive = req.Detail.IsMotherLive
		existingInvitation.InvitationDetail.IsFatherLive = req.Detail.IsFatherLive
		existingInvitation.InvitationDetail.IsBrideMotherLive = req.Detail.IsBrideMotherLive
		existingInvitation.InvitationDetail.IsBrideFatherLive = req.Detail.IsBrideFatherLive
		existingInvitation.InvitationDetail.IsGroomMotherLive = req.Detail.IsGroomMotherLive
		existingInvitation.InvitationDetail.IsGroomFatherLive = req.Detail.IsGroomFatherLive
	}

	if err := h.invitationService.UpdateInvitationWithRelations(c.UserContext(), existingInvitation); err != nil {
		if newFileName != "" {
			// DÜZELTME: Fonksiyonun dönüş değeri kullanılmadı.
			filemanager.DeleteFile("invitations", newFileName)
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye güncellenirken bir hata oluştu.")
		return c.Redirect(redirectURL, http.StatusSeeOther)
	}
	if oldPhotoToDelete != "" {
		// DÜZELTME: Fonksiyonun dönüş değeri kullanılmadı.
		filemanager.DeleteFile("invitations", oldPhotoToDelete)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla güncellendi.")
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}

func (h *DashboardInvitationHandler) DeleteInvitation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz davetiye ID'si.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	userID, _ := c.Locals("userID").(uint)
	ctxWithUser := context.WithValue(c.UserContext(), "user_id", userID)

	if err := h.invitationService.DeleteInvitationWithRelations(ctxWithUser, uint(id)); err != nil {
		errMsg := "Davetiye silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/invitations", fiber.StatusSeeOther)
	}
	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Davetiye başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla silindi.")
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}