package grpcserver

import (
	"context"

	purchasepb "github.com/CharLottery/proto/purchasepb"
	"github.com/CharLottery/purchase_service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PurchaseHandler struct {
	usecase *usecase.PurchaseUsecase
	purchasepb.UnimplementedPurchaseServiceServer
}

func NewPurchaseHandler(uc *usecase.PurchaseUsecase) *PurchaseHandler {
	return &PurchaseHandler{usecase: uc}
}

func (h *PurchaseHandler) BuyTicket(ctx context.Context, req *purchasepb.BuyTicketRequest) (*purchasepb.BuyTicketResponse, error) {
	if req.GetUserId() == 0 || req.GetLotteryId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id and lottery_id must be non-zero")
	}

	if len(req.GetNumbers()) != 5 {
		return nil, status.Error(codes.InvalidArgument, "exactly 5 numbers must be provided")
	}

	ticketID, err := h.usecase.BuyTicket(req.GetUserId(), req.GetLotteryId(), req.GetNumbers())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to buy ticket: %v", err)
	}

	return &purchasepb.BuyTicketResponse{
		Success:   true,
		Message:   "Ticket purchased successfully with your chosen numbers",
		TicketId:  ticketID,
		UserId:    req.GetUserId(),
		LotteryId: req.GetLotteryId(),
		Numbers:   req.GetNumbers(),
	}, nil
}
