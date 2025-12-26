package repository

import "github.com/jackc/pgx/v5/pgxpool"

type SubscriptionPlanRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionPlanRepository(db *pgxpool.Pool) *SubscriptionPlanRepository {
	return &SubscriptionPlanRepository{db: db}
}

func (r *SubscriptionPlanRepository) Create() {

}

func (r *SubscriptionPlanRepository) GetByID() {

}

func (r *SubscriptionPlanRepository) List() {

}

func (r *SubscriptionPlanRepository) Update() {

}
