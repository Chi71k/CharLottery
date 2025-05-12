package usecase

type LotteryEventPublisher interface {
	PublishLotteryCreated(lotteryID int64, prize string, availableTickets int64)
}
