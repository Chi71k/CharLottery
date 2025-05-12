package handler

func GetCardType(cardNumber string) string {
    if len(cardNumber) < 6 {
        return "Unknown"
    }

    prefix := cardNumber[:6]

    switch {
    case prefix == "440043":
        return "Kaspi Gold"
    case prefix == "404243":
        return "Forte Black"
    case prefix == "517792":
        return "Forte Blue"
    case prefix == "440563":
        return "Halyk Bonus"
    case prefix == "539545":
        return "Jusan Pay"
    case cardNumber[0:1] == "4":
        return "VISA"
    case cardNumber[0:2] >= "51" && cardNumber[0:2] <= "55":
        return "MASTERCARD"
    case cardNumber[0:2] == "34" || cardNumber[0:2] == "37":
        return "AMEX"
    default:
        return "Unknown"
    }
}
