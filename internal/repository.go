package user

import (
	"context"
	"fmt"
	"github.com/LuisRiveraBan/gocourse_domain/domain"
	"gorm.io/gorm"
	"log"
	"strings"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, firstname *string, lastName *string, email *string, phone *string) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repository struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database and returns an error if any occurred.
func (r *repository) Create(ctx context.Context, user *domain.User) error {

	// Add more fields as needed
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.log.Println(err)
		return err
	}

	// Log the successful creation of the user
	r.log.Println("user created with id: ", user.ID)
	return nil
}

// GetAll retrieves all users from the database and returns them as a slice of User structs.
func (r *repository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	var users []domain.User
	//return r.db.Find(&users).Error
	tx := r.db.WithContext(ctx).Model(&users)
	tx = applyFilters(tx, filters)
	tx = tx.Offset(offset).Limit(limit)
	if err := tx.Order("created_at desc").Find(&users).Error; err != nil {
		r.log.Println(err)
		return nil, err
	}
	return users, nil
}

// GetByID retrieves a user by their ID from the database and returns it as a User struct.
func (r *repository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	//return r.db.Where("id = ?", id).First(&user).Error
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		r.log.Println(err)
		return nil, ErrNotFound{id}
	}
	return &user, nil
}

// Delete deletes a user by their ID from the database.
func (r *repository) Delete(ctx context.Context, id string) error {
	user := domain.User{ID: id}

	resul := r.db.WithContext(ctx).Delete(&user)

	if resul.Error != nil {
		r.log.Println(resul.Error)
		return resul.Error
	}

	// Log the successful deletion of the user
	if resul.RowsAffected == 0 {
		r.log.Printf(
			"No user found with ID: %s", id,
		)
		return ErrNotFound{id}
	}
	// Log the successful deletion of the user
	return nil

}

func (r *repository) Update(ctx context.Context, id string, firstname *string, lastName *string, email *string, phone *string) error {
	values := make(map[string]interface{})

	if firstname != nil {
		values["first_name"] = *firstname
	}

	if lastName != nil {
		values["last_name"] = *lastName
	}

	if email != nil {
		values["email"] = *email
	}

	if phone != nil {
		values["phone"] = *phone
	}

	resul := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(values)

	if resul.Error != nil {
		r.log.Println(resul.Error)
		return resul.Error
	}

	if resul.RowsAffected == 0 {
		r.log.Printf(
			"No user found with ID: %s\n", id,
		)
		return ErrNotFound{id}
	}

	return nil
}

// con tx invocamos a la base de datos, y con filters a la estrutura de los filtros
func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	//si el FirstName es diferente que nil entonces con Tolwer lo volvera automaticamente el parametro
	//a minisculas para mayor porcentaje de exito
	//con tx ejecutamos un comando query para encontra informacion relacionada con el parametro
	//el %%%s%% es como decirle traeme todo lo que tenga eso
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}
	//
	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}

	return tx
}

// Count retrieves the total number of users matching the provided filters.
func (r *repository) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := r.db.Model(&domain.User{}).WithContext(ctx)
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		r.log.Println(err)
		return 0, err
	}
	return int(count), nil
}
