package repository

import (
	"github.com/med-000/notifyclass/pkg/logger"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewCourseRepository(db *gorm.DB, log *logger.RepositoryLogger) *CourseRepository {
	return &CourseRepository{
		db:  db,
		log: log,
	}
}

type ClassRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewClassRepository(db *gorm.DB, log *logger.RepositoryLogger) *ClassRepository {
	return &ClassRepository{
		db:  db,
		log: log,
	}
}

type GroupRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewGroupRepository(db *gorm.DB, log *logger.RepositoryLogger) *GroupRepository {
	return &GroupRepository{
		db:  db,
		log: log,
	}
}

type EventRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewEventRepository(db *gorm.DB, log *logger.RepositoryLogger) *EventRepository {
	return &EventRepository{
		db:  db,
		log: log,
	}
}

type ContentRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewContentRepository(db *gorm.DB, log *logger.RepositoryLogger) *ContentRepository {
	return &ContentRepository{
		db:  db,
		log: log,
	}
}

type NotionMappingRepository struct {
	db  *gorm.DB
	log *logger.RepositoryLogger
}

func NewNotionMappingRepository(db *gorm.DB, log *logger.RepositoryLogger) *NotionMappingRepository {
	return &NotionMappingRepository{
		db:  db,
		log: log,
	}
}
