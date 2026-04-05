package repository

import (
	"github.com/med-000/notifyclass/pkg/logger"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewCourseRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

type ClassRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

type GroupRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

type EventRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewEventRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

type ContentRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

type NotionMappingRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewNotionMappingRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}
