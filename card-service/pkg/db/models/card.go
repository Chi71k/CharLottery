package model

import (
    "errors"
    "regexp"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Card struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    CardNumber     string             `bson:"card_number"`
    CardHolder     string             `bson:"cardholder_name"`
    ExpiryDate     string             `bson:"expiration_date"`
    CVV            string             `bson:"cvv"`
    UserID         string             `bson:"user_id"`
    CardType       string             `bson:"card_type"`
}

func (c *Card) Validate() error {
    if len(c.CardNumber) < 12 {
        return errors.New("card number is too short")
    }

    if !regexp.MustCompile(`^(0[1-9]|1[0-2])\/\d{2}$`).MatchString(c.ExpiryDate) {
        return errors.New("invalid expiry date format, expected MM/YY")
    }

    // разбиваем на месяц и год
    parts := regexp.MustCompile(`/`).Split(c.ExpiryDate, -1)
    if len(parts) != 2 {
        return errors.New("invalid expiry date")
    }

    month := parts[0]
    year := parts[1]

    // добавляем префикс "20" к году (если YY < 50 — можно усложнить логику при необходимости)
    expiryTime, err := time.Parse("01/2006", month+"/20"+year)
    if err != nil {
        return errors.New("invalid expiry date")
    }

    // карта считается действительной до конца месяца
    endOfMonth := time.Date(expiryTime.Year(), expiryTime.Month()+1, 0, 23, 59, 59, 0, time.UTC)
    if time.Now().After(endOfMonth) {
        return errors.New("card is expired")
    }

    return nil
}

