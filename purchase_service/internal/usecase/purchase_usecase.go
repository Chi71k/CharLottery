package usecase

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "strconv"
  "time"

  "github.com/CharLottery/purchase_service/internal/adapter/postgres"
  "github.com/CharLottery/purchase_service/internal/model"
  "github.com/go-redis/redis/v8"
  "github.com/lib/pq"
  "gorm.io/gorm"
)

type PurchaseUsecase struct {
  repo      *postgres.PurchaseRepository
  publisher TicketEventPublisher
  cache     *redis.Client
  db        *gorm.DB
}

func NewPurchaseUsecase(repo *postgres.PurchaseRepository, publisher TicketEventPublisher, db *gorm.DB, cache *redis.Client) *PurchaseUsecase {
  return &PurchaseUsecase{repo: repo, publisher: publisher, db: db, cache: cache}
}

func (uc *PurchaseUsecase) BuyTicket(userIDStr string, lotteryID int64, numbers []int32) (int64, error) {
  userID, err := strconv.ParseInt(userIDStr, 10, 64)
  if err != nil {
    return 0, fmt.Errorf("invalid user_id format: %v", err)
  }

  numbersArray := pq.Int32Array(numbers)
  var ticketID int64

  err = uc.db.Transaction(func(tx *gorm.DB) error {
    var lotteryExists bool
    err := tx.Raw("SELECT EXISTS (SELECT 1 FROM lotteries WHERE id = ?)", lotteryID).Scan(&lotteryExists).Error
    if err != nil {
      return fmt.Errorf("failed to check if lottery exists: %v", err)
    }
    if !lotteryExists {
      return fmt.Errorf("lottery with ID %d does not exist", lotteryID)
    }

    var availableTickets int64
    err = tx.Raw("SELECT available_tickets FROM lotteries WHERE id = ?", lotteryID).Scan(&availableTickets).Error
    if err != nil {
      return fmt.Errorf("failed to check available tickets: %v", err)
    }
    if availableTickets <= 0 {
      return fmt.Errorf("no tickets available for this lottery")
    }

    var duplicatePurchase bool
    err = tx.Raw("SELECT EXISTS (SELECT 1 FROM purchases WHERE user_id = ? AND lottery_id = ? AND numbers = ?)",
      userID, lotteryID, numbersArray).Scan(&duplicatePurchase).Error
    if err != nil {
      return fmt.Errorf("failed to check for duplicate purchase: %v", err)
    }
    if duplicatePurchase {
      return fmt.Errorf("you have already purchased this combination of numbers for this lottery")
    }

    purchase := &model.Purchase{
      UserID:    userID,
      LotteryID: lotteryID,
      Numbers:   numbersArray,
    }
    if err := tx.Create(purchase).Error; err != nil {
      return fmt.Errorf("failed to create purchase: %v", err)
    }
    ticketID = purchase.ID

    err = tx.Exec("UPDATE lotteries SET available_tickets = available_tickets - 1 WHERE id = ?", lotteryID).Error
    if err != nil {
      return fmt.Errorf("failed to update available tickets: %v", err)
    }

    uc.cache.Del(context.Background(), fmt.Sprintf("user_tickets_%d", userID))

    uc.publisher.PublishTicketBought(lotteryID)

    return nil
  })

  if err != nil {
    log.Printf("Transaction failed: %v", err)
    return 0, err
  }

  log.Println("Transaction succeeded, ticket purchased successfully")
  return ticketID, nil
}

func (uc *PurchaseUsecase) ListTicketsByUser(userIDStr string) ([]model.Purchase, error) {
  userID, err := strconv.ParseInt(userIDStr, 10, 64)
  if err != nil {
    return nil, fmt.Errorf("invalid user_id format: %v", err)
  }

  cacheKey := fmt.Sprintf("user_tickets_%d", userID)
  ctx := context.Background()

  log.Printf("Attempting to get tickets for user %d from cache with key %s", userID, cacheKey)

  val, err := uc.cache.Get(ctx, cacheKey).Result()
  if err == nil {
    var cached []model.Purchase
    if err := json.Unmarshal([]byte(val), &cached); err == nil {
      log.Printf("Cache hit: tickets found for user %d", userID)
      return cached, nil
    } else {
      log.Printf("Cache hit, but failed to unmarshal data for user %d: %v", userID, err)
    }
  }

  log.Printf("Cache miss for user %d. Fetching tickets from database.", userID)

  tickets, err := uc.repo.ListByUser(userID)
  if err != nil {
    log.Printf("Failed to fetch tickets for user %d from database: %v", userID, err)
    return nil, err
  }


  log.Printf("Successfully fetched tickets for user %d from database. Saving to cache.", userID)

  data, err := json.Marshal(tickets)
  if err != nil {
    log.Printf("Failed to marshal tickets for user %d: %v", userID, err)
  } else {
    uc.cache.Set(ctx, cacheKey, data, 5*time.Minute)
    log.Printf("Successfully saved tickets for user %d to cache", userID)
  }

  return tickets, nil
}

func (uc *PurchaseUsecase) UpdatePurchase(purchaseID int64, userIDStr string, newNumbers []int32) error {
  userID, err := strconv.ParseInt(userIDStr, 10, 64)
  if err != nil {
    return fmt.Errorf("invalid user_id format: %v", err)
  }

  numbersArray := pq.Int32Array(newNumbers)

  err = uc.db.Transaction(func(tx *gorm.DB) error {
    var exists bool
    err := tx.Raw("SELECT EXISTS (SELECT 1 FROM purchases WHERE id = ? AND user_id = ?)", purchaseID, userID).Scan(&exists).Error
    if err != nil {
      return fmt.Errorf("failed to check if purchase exists: %v", err)
    }
    if !exists {
      return fmt.Errorf("purchase not found or does not belong to user")
    }

    err = tx.Exec("UPDATE purchases SET numbers = ? WHERE id = ? AND user_id = ?", numbersArray, purchaseID, userID).Error
    if err != nil {
      return fmt.Errorf("failed to update purchase: %v", err)
    }

    cacheKey := fmt.Sprintf("user_tickets_%d", userID)
    uc.cache.Del(context.Background(), cacheKey)

    return nil
  })

  if err != nil {
    log.Printf("Failed to update purchase %d for user %d: %v", purchaseID, userID, err)
    return err
  }

  log.Printf("Purchase %d for user %d updated successfully", purchaseID, userID)
  return nil
}

func (uc *PurchaseUsecase) DeletePurchase(purchaseID int64, userIDStr string) error {
  userID, err := strconv.ParseInt(userIDStr, 10, 64)
  if err != nil {
    return fmt.Errorf("invalid user_id format: %v", err)
  }

  err = uc.db.Transaction(func(tx *gorm.DB) error {
    var exists bool
    err := tx.Raw("SELECT EXISTS (SELECT 1 FROM purchases WHERE id = ? AND user_id = ?)", purchaseID, userID).Scan(&exists).Error
    if err != nil {
      return fmt.Errorf("failed to check if purchase exists: %v", err)
    }
    if !exists {
      return fmt.Errorf("purchase not found or does not belong to user")
    }

    err = tx.Exec("DELETE FROM purchases WHERE id = ? AND user_id = ?", purchaseID, userID).Error
    if err != nil {
      return fmt.Errorf("failed to delete purchase: %v", err)
    }

    cacheKey := fmt.Sprintf("user_tickets_%d", userID)
    uc.cache.Del(context.Background(), cacheKey)

    return nil
  })

  if err != nil {
    log.Printf("Failed to delete purchase %d for user %d: %v", purchaseID, userID, err)
    return err
  }

  log.Printf("Purchase %d for user %d deleted successfully", purchaseID, userID)
  return nil
}

type TicketEventPublisher interface {
  PublishTicketBought(lotteryID int64)
}
