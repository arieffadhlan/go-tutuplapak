package entities

import "time"

type (
	Merchant struct {
		MerchantID       string           `db:"merchant_id"`
		Name             string           `db:"name"`
		MerchantCategory MerchantCategory `db:"merchant_category"`
		ImageURL         string           `db:"image_url"`
		Location         Location         `db:"location"`
		CreatedAt        time.Time        `db:"created_at"`
	}

	// MerchantCategory enum type in Go (aligns with Postgres ENUM)
	MerchantCategory string

	// Location represents PostGIS Point (4326)
	Location struct {
		Lat  float64
		Long float64
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
