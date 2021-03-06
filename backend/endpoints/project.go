package endpoints

import (
	"time"

	"github.com/1backend/1backend/backend/deploy"
	"github.com/1backend/1backend/backend/domain"
	tp "github.com/1backend/1backend/backend/tech-pack"
)

func (e Endpoints) CreateProject(proj *domain.Project) error {
	proj.Description = "Change this to something sensible"
	proj.Id = domain.Sid.MustGenerate()
	proj.InfraPassword = domain.Sid.MustGenerate()
	proj.CreatedAt = time.Now()
	proj.UpdatedAt = time.Now()
	tp.GetProvider(proj).CreateProjectPlugin()
	err := e.db.Save(proj).Error
	if err != nil {
		return err
	}
	for _, v := range proj.Endpoints {
		v.Id = domain.Sid.MustGenerate()
		v.ProjectId = proj.Id
		v.CreatedAt = time.Now()
		v.UpdatedAt = time.Now()
		err = e.db.Save(&v).Error
		if err != nil {
			return err
		}
	}
	for _, v := range proj.Dependencies {
		v.Id = domain.Sid.MustGenerate()
		v.ProjectId = proj.Id
		v.CreatedAt = time.Now()
		v.UpdatedAt = time.Now()
		err = e.db.Save(&v).Error
		if err != nil {
			return err
		}
	}
	go deploy.NewDeployer(e.db).Deploy(proj)
	return nil
}

func (e Endpoints) UpdateProject(proj *domain.Project) error {
	proj.InfraPassword = domain.Sid.MustGenerate()
	err := e.db.Save(proj).Error
	if err != nil {
		return err
	}
	for _, v := range proj.Endpoints {
		if v.Id == "" {
			v.Id = domain.Sid.MustGenerate()
		}
		v.ProjectId = proj.Id
		err = e.db.Save(&v).Error
		if err != nil {
			return err
		}
	}
	for _, v := range proj.Dependencies {
		v.Id = domain.Sid.MustGenerate()
		v.ProjectId = proj.Id
		err = e.db.Save(&v).Error
		if err != nil {
			return err
		}
	}
	go deploy.NewDeployer(e.db).Deploy(proj)
	return nil
}

func (e Endpoints) PutStar(userId, projectId string) error {
	star := &domain.Star{
		Id:        domain.Sid.MustGenerate(),
		ProjectId: projectId,
		UserId:    userId,
	}
	return e.db.Save(star).Error
}

func (e Endpoints) DeleteStar(userId, projectId string) error {
	star := domain.Star{}
	err := e.db.Where("user_id = ? AND project_id = ?", userId, projectId).First(&star).Error
	if err != nil {
		return err
	}
	return e.db.Delete(star).Error
}
