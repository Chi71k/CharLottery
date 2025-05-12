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

    exp, err := time.Parse("01/06", c.ExpiryDate)
    if err != nil || exp.Before(time.Now()) {
        return errors.New("card is expired")
    }

    return nil
}
