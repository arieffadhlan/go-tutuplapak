package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tutuplapak-user/internal/dto"
	"tutuplapak-user/internal/utils"

	"github.com/jmoiron/sqlx"
)

type MerchantsBeliRepositoryInterface interface {
	CreateMerchant(ctx context.Context, req dto.MerchantCreateRequest) (merchant dto.MerchantCreateResponse, err error)
	GetMerchant(ctx context.Context, filter dto.MerchantFilter) (list dto.ListMerchantResponse, err error)
	GetMerchantById(ctx context.Context, merchantId string) (merchant dto.Merchant, err error)
	CreateItem(ctx context.Context, merchantId string, req dto.ItemCreateRequest) (item dto.ItemCreateResponse, err error)
	GetItem(ctx context.Context, merchantId string, filter dto.ItemFilter) (list dto.ListItemResponse, err error)
}

type MerchantsBeliRepository struct {
	db *sqlx.DB
}

func NewMerchantsBeliRepository(db *sqlx.DB) MerchantsBeliRepository {
	return MerchantsBeliRepository{db: db}
}

func (r MerchantsBeliRepository) CreateMerchant(ctx context.Context, req dto.MerchantCreateRequest) (merchant dto.MerchantCreateResponse, err error) {
	query := `
		INSERT INTO merchants (name, merchant_category, image_url, location)
		VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326)::GEOGRAPHY)
		RETURNING merchant_id, name, merchant_category, image_url,
		          ST_Y(location::geometry) AS lat,
		          ST_X(location::geometry) AS long,
		          created_at
	`
	err = r.db.QueryRowContext(ctx, query, req.Name, req.MerchantCategory, req.ImageURL, req.Location.Long, req.Location.Lat).Scan(&merchant.ID)

	if err != nil {
		return dto.MerchantCreateResponse{}, err
	}

	return merchant, nil
}

func (r MerchantsBeliRepository) GetMerchantById(ctx context.Context, merchantId string) (merchant dto.Merchant, err error) {
	query := `
		SELECT 
			merchant_id,
			name,
			merchant_category,
			image_url,
			ST_Y(location::geometry) as lat,
			ST_X(location::geometry) as long,
			created_at
		FROM merchants
		WHERE merchant_id = $1
	`

	err = r.db.GetContext(ctx, &merchant, query, merchantId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Merchant{}, utils.NewNotFound("file not found")
		}
		return dto.Merchant{}, utils.NewInternal("failed to get file")
	}

	return merchant, nil
}

func (r MerchantsBeliRepository) GetMerchant(ctx context.Context, filter dto.MerchantFilter) (list dto.ListMerchantResponse, err error) {
	args := []interface{}{}
	conditions := []string{"1=1"} // always true, so filters can be appended

	if filter.MerchantID != "" {
		conditions = append(conditions, "merchant_id = ?")
		args = append(args, filter.MerchantID)
	}

	if filter.Name != "" {
		conditions = append(conditions, "LOWER(name) LIKE ?")
		args = append(args, "%"+strings.ToLower(filter.Name)+"%")
	}

	if filter.MerchantCategory != "" {
		conditions = append(conditions, "merchant_category = ?")
		args = append(args, filter.MerchantCategory)
	}

	// Base query
	baseQuery := `
		SELECT 
			merchant_id,
			name,
			merchant_category,
			image_url,
			ST_Y(location::geometry) as lat,
			ST_X(location::geometry) as long,
			created_at
		FROM merchants
		WHERE ` + strings.Join(conditions, " AND ")

	// Sorting
	if filter.SortCreatedAt == "asc" || filter.SortCreatedAt == "desc" {
		baseQuery += fmt.Sprintf(" ORDER BY created_at %s", filter.SortCreatedAt)
	}

	// Pagination
	limit := 5
	offset := 0
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Offset >= 0 {
		offset = filter.Offset
	}
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Run query
	var merchants []dto.Merchant
	err = r.db.SelectContext(ctx, &merchants, r.db.Rebind(baseQuery), args...)
	if err != nil {
		return dto.ListMerchantResponse{}, err
	}

	// Count query for meta
	countQuery := `
		SELECT COUNT(*) 
		FROM merchants
		WHERE ` + strings.Join(conditions, " AND ")
	var total int
	err = r.db.GetContext(ctx, &total, r.db.Rebind(countQuery), args...)
	if err != nil {
		return dto.ListMerchantResponse{}, err
	}

	return dto.ListMerchantResponse{
		Data: merchants,
		Meta: dto.Meta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}, nil
}

func (r MerchantsBeliRepository) CreateItem(ctx context.Context, merchantId string, req dto.ItemCreateRequest) (item dto.ItemCreateResponse, err error) {
	query := `
		INSERT INTO items (merchant_id, name, product_category, price, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING item_id, merchant_id, name, product_category, price, image_url, created_at
	`
	err = r.db.QueryRowContext(ctx, query, merchantId, req.Name, req.ProductCategory, req.Price, req.ImageURL).Scan(&item.ID)
	if err != nil {
		return dto.ItemCreateResponse{}, err
	}

	return item, nil
}

func (r MerchantsBeliRepository) GetItem(ctx context.Context, merchantId string, filter dto.ItemFilter) (list dto.ListItemResponse, err error) {
	args := []interface{}{merchantId}
	conditions := []string{"merchant_id = ?"} // required

	if filter.ItemID != "" {
		conditions = append(conditions, "item_id = ?")
		args = append(args, filter.ItemID)
	}

	if filter.Name != "" {
		conditions = append(conditions, "LOWER(name) LIKE ?")
		args = append(args, "%"+strings.ToLower(filter.Name)+"%")
	}

	if filter.ProductCategory != "" {
		conditions = append(conditions, "product_category = ?")
		args = append(args, filter.ProductCategory)
	}

	// Base query
	baseQuery := `
		SELECT 
			item_id,
			name,
			product_category,
			price,
			image_url,
			created_at
		FROM items
		WHERE ` + strings.Join(conditions, " AND ")

	// Sorting
	if filter.SortCreatedAt == "asc" || filter.SortCreatedAt == "desc" {
		baseQuery += fmt.Sprintf(" ORDER BY created_at %s", filter.SortCreatedAt)
	}

	// Pagination
	limit := 5
	offset := 0
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Offset >= 0 {
		offset = filter.Offset
	}
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Query rows
	var items []dto.Item
	err = r.db.SelectContext(ctx, &items, r.db.Rebind(baseQuery), args...)
	if err != nil {
		return dto.ListItemResponse{}, err
	}

	// Count for meta
	countQuery := `
		SELECT COUNT(*)
		FROM items
		WHERE ` + strings.Join(conditions, " AND ")
	var total int
	err = r.db.GetContext(ctx, &total, r.db.Rebind(countQuery), args...)
	if err != nil {
		return dto.ListItemResponse{}, err
	}

	return dto.ListItemResponse{
		Data: items,
		Meta: dto.Meta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}, nil
}
