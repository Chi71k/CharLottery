package grpcserver

import (
	"context"
	"fmt"

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
	userIDFromCtx, ok := ctx.Value("userID").(string)
	if !ok || userIDFromCtx == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	reqUserID := fmt.Sprintf("%d", req.GetUserId())
	if reqUserID != userIDFromCtx {
		return nil, status.Error(codes.PermissionDenied, "user ID does not match authenticated user")
	}

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

func (h *PurchaseHandler) ListTicketsByUser(ctx context.Context, req *purchasepb.ListTicketsByUserRequest) (*purchasepb.ListTicketsByUserResponse, error) {
	userIDFromCtx, ok := ctx.Value("userID").(string)
	if !ok || userIDFromCtx == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	reqUserID := fmt.Sprintf("%d", req.GetUserId())
	if reqUserID != userIDFromCtx {
		return nil, status.Error(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id must be non-zero")
	}

	tickets, err := h.usecase.ListTicketsByUser(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tickets: %v", err)
	}

	var pbTickets []*purchasepb.Ticket
	for _, t := range tickets {
		pbTickets = append(pbTickets, &purchasepb.Ticket{
			TicketId:  t.ID,
			UserId:    t.UserID,
			LotteryId: t.LotteryID,
			Numbers:   t.Numbers,
		})
	}

	return &purchasepb.ListTicketsByUserResponse{Tickets: pbTickets}, nil
}

func (h *PurchaseHandler) UpdatePurchase(ctx context.Context, req *purchasepb.UpdatePurchaseRequest) (*purchasepb.UpdatePurchaseResponse, error) {
	userIDFromCtx, ok := ctx.Value("userID").(string)
	if !ok || userIDFromCtx == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	reqUserID := fmt.Sprintf("%d", req.GetUserId())
	if reqUserID != userIDFromCtx {
		return nil, status.Error(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	if req.GetPurchaseId() == 0 || req.GetUserId() == 0 || len(req.GetNewNumbers()) != 5 {
		return nil, status.Error(codes.InvalidArgument, "purchase_id, user_id and exactly 5 new numbers are required")
	}

	err := h.usecase.UpdatePurchase(req.GetPurchaseId(), req.GetUserId(), req.GetNewNumbers())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update purchase: %v", err)
	}

	return &purchasepb.UpdatePurchaseResponse{
		Success: true,
		Message: "Purchase updated successfully",
	}, nil
}

func (h *PurchaseHandler) DeletePurchase(ctx context.Context, req *purchasepb.DeletePurchaseRequest) (*purchasepb.DeletePurchaseResponse, error) {
	userIDFromCtx, ok := ctx.Value("userID").(string)
	if !ok || userIDFromCtx == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}
	reqUserID := fmt.Sprintf("%d", req.GetUserId())
	if reqUserID != userIDFromCtx {
		return nil, status.Error(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	if req.GetPurchaseId() == 0 || req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "purchase_id and user_id are required")
	}

	err := h.usecase.DeletePurchase(req.GetPurchaseId(), req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete purchase: %v", err)
	}

	return &purchasepb.DeletePurchaseResponse{
		Success: true,
		Message: "Purchase deleted successfully",
	}, nil
}
