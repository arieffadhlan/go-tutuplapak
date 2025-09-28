package dto

import "time"

type (
	Merchant struct {
		MerchantID       string           `json:"merchantId" db:"merchant_id"`
		Name             string           `json:"name" db:"name"`
		MerchantCategory MerchantCategory `json:"merchantCategory" db:"merchant_category"`
		ImageURL         string           `json:"imageUrl" db:"image_url"`
		Location         Location         `json:"location" db:"location"`
		CreatedAt        time.Time        `json:"createdAt" db:"created_at"`
	}

	MerchantCategory string

	Location struct {
		Lat  float64 `json:"lat" validate:"required"`
		Long float64 `json:"long" validate:"required"`
	}

	Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	MerchantCreateResponse struct {
		ID string `json:"merchantId"`
	}

	// Response wrapper for API
	ListMerchantResponse struct {
		Data []Merchant `json:"data"`
		Meta Meta       `json:"meta"`
	}

	MerchantCreateRequest struct {
		Name             string           `json:"name" validate:"required,min=2,max=30"`
		MerchantCategory MerchantCategory `json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
		ImageURL         string           `json:"imageUrl" validate:"required"`
		Location         Location         `json:"location"`
	}

	MerchantFilter struct {
		MerchantID       string
		Name             string
		MerchantCategory string
		Limit            int
		Offset           int
		SortCreatedAt    string // "asc" or "desc"
	}

	Item struct {
		ItemId          string          `json:"itemId" db:"item_id"`
		Name            string          `json:"name" db:"name"`
		ProductCategory ProductCategory `json:"productCategory" db:"product_category"`
		Price           int             `json:"price" db:"price"`
		ImageURL        string          `json:"imageUrl" db:"image_url"`
		CreateAt        time.Time       `json:"createdAt" db:"created_at"`
	}

	ItemCreateRequest struct {
		Name            string          `json:"name" validate:"required,min=2,max=30"`
		ProductCategory ProductCategory `json:"productCategory" validate:"required, oneof=Beverage Food Snack Condiments Additions"`
		Price           int             `json:"price" validate:"required,min=1"`
		ImageURL        string          `json:"imageUrl" validate:"required"`
	}
	ProductCategory string

	ItemFilter struct {
		ItemID          string
		Limit           int
		Offset          int
		Name            string
		ProductCategory string
		SortCreatedAt   string // "asc" or "desc"
	}

	ItemCreateResponse struct {
		ID string `json:"itemId"`
	}
	ListItemResponse struct {
		Data []Item `json:"data"`
		Meta Meta   `json:"meta"`
	}
)

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

const (
	Beverage   ProductCategory = "Beverage"
	Food       ProductCategory = "Food"
	Snack      ProductCategory = "Snack"
	Condiments ProductCategory = "Condiments"
	Additions  ProductCategory = "Additions"
)
