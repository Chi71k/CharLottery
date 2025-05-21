package grpcserver

import (
  "context"
  "errors"
  "fmt"
  "strings"
  "strconv"

  purchasepb "github.com/CharLottery/proto/purchasepb"
  "github.com/CharLottery/purchase_service/internal/usecase"

  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"

  "github.com/golang-jwt/jwt/v5"
)

// JWT секретный ключ
var secretKey = []byte("maxsecretkey")

// AuthUnaryInterceptor — middleware для проверки JWT
func AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
  return func(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
  ) (interface{}, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
      return nil, errors.New("missing metadata")
    }

    authHeader := md.Get("authorization")
    if len(authHeader) == 0 {
      return nil, errors.New("missing authorization header")
    }

    tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")
    userID, err := validateToken(tokenStr)
    if err != nil {
      return nil, err
    }

    // Добавляем userID в контекст
    ctx = context.WithValue(ctx, "userID", userID)
    return handler(ctx, req)
  }
}

// validateToken парсит и проверяет JWT, возвращает userID из claims
func validateToken(tokenStr string) (string, error) {
  token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return secretKey, nil
  })

  if err != nil {
    return "", err
  }

  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    userID, ok := claims["user_id"].(string)
    if !ok {
      return "", errors.New("user_id not found in token")
    }
    return userID, nil
  }

  return "", errors.New("invalid token")
}

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
  if req.GetUserId() != userIDFromCtx {
    return nil, status.Error(codes.PermissionDenied, "user ID does not match authenticated user")
  }

  if req.GetUserId() == "" || req.GetLotteryId() == 0 {
    return nil, status.Error(codes.InvalidArgument, "user_id and lottery_id must be provided")
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
  if req.GetUserId() != userIDFromCtx {
    return nil, status.Error(codes.PermissionDenied, "user ID does not match authenticated user")
  }

  if req.GetUserId() == "" {
    return nil, status.Error(codes.InvalidArgument, "user_id must be provided")
  }

  tickets, err := h.usecase.ListTicketsByUser(req.GetUserId())
  if err != nil {
    return nil, status.Errorf(codes.Internal, "failed to list tickets: %v", err)
  }


  var pbTickets []*purchasepb.Ticket
  for _, t := range tickets {
    pbTickets = append(pbTickets, &purchasepb.Ticket{
      TicketId:  t.ID,
      UserId:    strconv.FormatInt(t.UserID, 10),
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
  if req.GetUserId() != userIDFromCtx {
    return nil, status.Error(codes.PermissionDenied, "user ID does not match authenticated user")
  }

  if req.GetPurchaseId() == 0 || req.GetUserId() == "" || len(req.GetNewNumbers()) != 5 {
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
  if req.GetUserId() != userIDFromCtx {
    return nil, status.Error(codes.PermissionDenied, "user ID does not match authenticated user")
  }

  if req.GetPurchaseId() == 0 || req.GetUserId() == "" {
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
