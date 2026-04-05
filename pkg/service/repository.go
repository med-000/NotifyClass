package service

import (
	"github.com/med-000/notifyclass/pkg/logger"
	"github.com/med-000/notifyclass/pkg/parser"
	"github.com/med-000/notifyclass/pkg/repository"
	"gorm.io/gorm"
)

func (s *Service) SaveAll(dbConn *gorm.DB, course *parser.Course) error {
	repositoryLogger, _ := logger.NewRepositoryLogger()
	courseRepo := repository.NewCourseRepository(dbConn,repositoryLogger)
	classRepo := repository.NewClassRepository(dbConn,repositoryLogger)
	groupRepo := repository.NewGroupRepository(dbConn,repositoryLogger)
	eventRepo := repository.NewEventRepository(dbConn,repositoryLogger)
	// contentRepo := repository.NewContentRepository(dbConn)

	err := courseRepo.Save(course)
	if err != nil {
		s.log.Error.Printf("defined course save err:%s", err)
		return err
	}
	s.log.Info.Printf("save course")

	for _, class := range course.Classes {
		err := classRepo.Save(class)
		if err != nil {
			s.log.Error.Printf("defined class save err:%s", err)
			return err
		}
		for _, group := range class.Groups {
			err := groupRepo.Save(group)
			if err != nil {
				s.log.Error.Printf("defined group save err:%s", err)
				return err
			}
			for _, event := range group.Events {
				err := eventRepo.Save(event)
				if err != nil {
					s.log.Error.Printf("defined event save err:%s", err)
					return err
				}
				// for _, content := range event.Content {
				// 	err := contentRepo.Save(content)
				// 	if err != nil {
				// 		s.log.Error.Printf("defined content save err:%s", err)
				// 		return err
				// 	}
				// }
				// s.log.Info.Printf("save content")
			}
			s.log.Info.Printf("save event")
		}
		s.log.Info.Printf("save group")
	}
	s.log.Info.Printf("save class")
	return nil
}
