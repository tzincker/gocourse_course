package course

import (
	"context"
	"log"
	"time"

	"github.com/tzincker/gocourse_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, name *string, startDate *string, endDate *string) error
		Count(ctx context.Context, filters Filters) (int64, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error) {
	log.Println("Create course service")

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Panicln(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Panicln(err)
		return nil, err
	}

	course := domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	u, err := s.repo.Create(ctx, &course)

	if err != nil {
		s.log.Println(err)
	}

	return u, err
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	log.Println("Get all courses service")

	courses, err := s.repo.GetAll(ctx, filters, offset, limit)

	if err != nil {
		s.log.Println(err)
	}

	return courses, err
}

func (s service) Get(ctx context.Context, id string) (*domain.Course, error) {
	log.Println("Get course service")

	course, err := s.repo.Get(ctx, id)

	if err != nil {
		s.log.Println(err)
	}

	return course, err
}

func (s service) Delete(ctx context.Context, id string) error {
	log.Println("Delete course service")

	err := s.repo.Delete(ctx, id)

	if err != nil {
		s.log.Println(err)
		return err
	}

	return nil
}

func (s service) Update(
	ctx context.Context,
	id string,
	name *string,
	startDate *string,
	endDate *string,
) error {
	log.Println("Update course service")

	var err error
	var startDateParsed *time.Time
	var endDateParsed *time.Time

	if startDate != nil {

		parsedDate, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Panicln(err)
			return err
		}
		startDateParsed = &parsedDate
	}

	if endDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Panicln(err)
			return err
		}
		endDateParsed = &parsedDate
	}

	err = s.repo.Update(ctx, id, name, startDateParsed, endDateParsed)
	if err != nil {
		s.log.Println(err)
		return err
	}

	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int64, error) {
	log.Println("Get all courses count service")
	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		s.log.Println(err)
	}

	return count, err
}
