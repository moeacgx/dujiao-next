package service

import (
	"fmt"
	"strings"

	"github.com/dujiao-next/internal/models"
	"github.com/dujiao-next/internal/repository"
)

// MediaService 素材管理服务
type MediaService struct {
	repo repository.MediaRepository
}

// NewMediaService 创建素材服务实例
func NewMediaService(repo repository.MediaRepository) *MediaService {
	return &MediaService{repo: repo}
}

// List 素材列表
func (s *MediaService) List(scene, search string, page, pageSize int) ([]models.Media, int64, error) {
	return s.repo.List(repository.MediaListFilter{
		Page:     page,
		PageSize: pageSize,
		Scene:    scene,
		Search:   search,
	})
}

// RecordMedia 记录上传的素材元数据（上传后自动调用）
func (s *MediaService) RecordMedia(result *UploadResult, scene string) (*models.Media, error) {
	// 检查是否已存在（基于路径去重）
	existing, err := s.repo.GetByPath(result.URL)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	// 从原始文件名生成默认素材名称（去掉扩展名）
	name := result.Filename
	if idx := strings.LastIndex(name, "."); idx > 0 {
		name = name[:idx]
	}

	media := &models.Media{
		Name:     name,
		Filename: result.Filename,
		Path:     result.URL,
		MimeType: result.MimeType,
		Size:     result.Size,
		Scene:    scene,
		Width:    result.Width,
		Height:   result.Height,
	}
	if err := s.repo.Create(media); err != nil {
		return nil, err
	}
	return media, nil
}

// Rename 重命名素材
func (s *MediaService) Rename(id uint, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("素材名称不能为空")
	}
	media, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if media == nil {
		return fmt.Errorf("素材不存在")
	}
	media.Name = name
	return s.repo.Update(media)
}

// Delete 删除素材（仅软删除记录，不删物理文件）
func (s *MediaService) Delete(id uint) error {
	media, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if media == nil {
		return fmt.Errorf("素材不存在")
	}
	return s.repo.Delete(id)
}
