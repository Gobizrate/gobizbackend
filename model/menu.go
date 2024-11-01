package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	Name          string `bson:"name" json:"name"`
	Price         int    `bson:"price" json:"price"`
	OriginalPrice int    `bson:"originalPrice" json:"originalPrice"`
	Rating        Rating `bson:"rating" json:"rating"`
	Sold          int    `bson:"sold" json:"sold"`
	Image         string `bson:"image" json:"image"`
}

type Toko struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NamaToko   string             `bson:"nama_toko" json:"nama_toko"`
	Slug       string             `bson:"slug" json:"slug"`
	Category   Category           `bson:"category" json:"category"`
	GambarToko string             `bson:"gambar_toko" json:"gambar_toko"`
	Alamat     Address            `bson:"alamat" json:"alamat"`
	User       []Userdomyikado    `bson:"user" json:"user"`
	Menu       []Menu             `bson:"menu" json:"menu"`
	Rating     Rating             `bson:"rating" json:"rating"`
}

type Address struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Street      string             `bson:"street" json:"street,omitempty"`
	Province    string             `bson:"province" json:"province,omitempty"`
	City        string             `bson:"city" json:"city,omitempty"`
	Description string             `bson:"description" json:"description,omitempty"`
	PostalCode  string             `bson:"postal_code" json:"postal_code,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type Category struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CategoryName string             `bson:"name_category" json:"name_category"`
}

type Rating struct {
	Average float64 `bson:"average" json:"average"`
	Count   int     `bson:"count" json:"count"`
}

type FAQ struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nama      string             `bson:"nama" json:"nama,omitempty"`
	Email     string             `bson:"email" json:"email,omitempty"`
	Handphone string             `bson:"handphone" json:"handphone,omitempty"`
	Questions string             `bson:"questions" json:"questions,omitempty"`
}
